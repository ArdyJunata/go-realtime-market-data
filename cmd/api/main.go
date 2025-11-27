package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"github.com/ArdyJunata/go-realtime-market-data/internal/handler.go"
	"github.com/ArdyJunata/go-realtime-market-data/internal/repository"
	"github.com/ArdyJunata/go-realtime-market-data/internal/service"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/database"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.LoadConfig("../../.env")
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
	_ = handler.NewHandler(service)

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
		api.Get("/health", func(c *fiber.Ctx) error {
			fmt.Println("OK")
			return nil
		})
	}

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
