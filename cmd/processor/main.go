package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"github.com/ArdyJunata/go-realtime-market-data/internal/repository"
	"github.com/ArdyJunata/go-realtime-market-data/internal/service"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/constant"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/database"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/logger"
)

func main() {
	config.LoadConfig(".env")

	logger.InitLogger()
	ctx := context.Background()

	logger.Log.Infof(ctx, "Starting Market Processor...")

	rdb, err := database.InitRedis()
	if err != nil {
		logger.Log.Errorf(ctx, "Redis Error: %v", err)
		os.Exit(1)
	}

	mongoClient, err := database.InitMongo()
	if err != nil {
		logger.Log.Errorf(ctx, "Mongo Error: %v", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(mongoClient, rdb)
	service := service.NewService(repo)

	logger.Log.Infof(ctx, "Subscribing to channel: %s", constant.RedisChannelMarketTrades)
	pubsub := rdb.Subscribe(ctx, constant.RedisChannelMarketTrades)

	defer pubsub.Close()

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			err := service.ProcessTradeEvent(ctx, msg.Payload)

			if err != nil {
				logger.Log.Warnf(ctx, "Invalid message format: %v", err)
			}

		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Warnf(ctx, "Shutting down processor...")
}
