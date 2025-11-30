package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"github.com/ArdyJunata/go-realtime-market-data/internal/handler"
	"github.com/ArdyJunata/go-realtime-market-data/internal/repository"
	"github.com/ArdyJunata/go-realtime-market-data/internal/service"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/database"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/logger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.LoadConfig(".env")

	logger.InitLogger()
	ctx := context.Background()

	logger.Log.Infof(ctx, "Starting Market API Service with Fiber...")

	mongo, err := database.InitMongo()
	if err != nil {
		logger.Log.Errorf(ctx, err.Error())
		panic(err)
	}

	redis, err := database.InitRedis()
	if err != nil {
		logger.Log.Errorf(ctx, err.Error())
		panic(err)
	}

	repo := repository.NewRepository(mongo, redis)
	service := service.NewService(repo)
	hub := handler.NewHub()
	handler := handler.NewHandler(service, hub)

	go func() {
		logger.Log.Infof(ctx, "Starting Data Pipeline: Redis -> Service -> Hub")
		streamChan := service.GetTradeStream(ctx)
		for msg := range streamChan {
			hub.Broadcast([]byte(msg))
		}
	}()

	router := fiber.New(fiber.Config{
		AppName:        "Realtime Market Data",
		RequestMethods: fiber.DefaultMethods,
	})

	router.Use(recover.New())
	router.Use(cors.New())
	router.Use(fiberlogger.New(fiberlogger.Config{
		TimeZone: "Asia/Jakarta",
	}))

	api := router.Group("/api/v1")
	{
		api.Get("/price/:symbol", handler.GetPriceSnapshot)
		api.Get("/trades/:symbol", handler.GetTrades)
	}

	router.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	router.Get("/ws", websocket.New(handler.HandleWebSocket))

	go func() {
		logger.Log.Infof(ctx, "Server listening on port 8080")
		if err := router.Listen(":8080"); err != nil {
			logger.Log.Errorf(ctx, "Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Warnf(ctx, "Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := router.Shutdown(); err != nil {
		logger.Log.Errorf(ctx, "Server forced to shutdown: %v", err)
	}

	if err := mongo.Disconnect(ctx); err != nil {
		logger.Log.Errorf(ctx, "Error disconnecting Mongo: %v", err)
	}

	logger.Log.Infof(ctx, "Server exiting")
}
