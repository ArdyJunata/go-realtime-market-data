package database

import (
	"context"
	"fmt"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetString(config.CFG_REDIS_HOST), config.GetString(config.CFG_REDIS_PORT)),
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return rdb, nil
}
