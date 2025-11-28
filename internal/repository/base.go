package repository

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct {
	mongo *mongo.Client
	redis *redis.Client
}

func NewRepository(mongoDb *mongo.Client, redis *redis.Client) *repository {
	return &repository{
		mongo: mongoDb,
		redis: redis,
	}
}
