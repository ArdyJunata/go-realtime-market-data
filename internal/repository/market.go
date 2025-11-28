package repository

import (
	"context"
	"fmt"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/constant"
)

func (r *repository) InsertTrade(ctx context.Context, trade entity.Trade) error {
	collection := r.mongo.Database(constant.MongoDatabaseName).Collection(constant.MongoCollectionTrades)

	_, err := collection.InsertOne(ctx, trade)
	return err
}

func (r *repository) UpdateLastPrice(ctx context.Context, trade entity.Trade) error {

	key := fmt.Sprintf("%s%s", constant.RedisKeyLatestPrice, trade.Symbol)

	return r.redis.Set(ctx, key, trade.Price, 0).Err()
}
