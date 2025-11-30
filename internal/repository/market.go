package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/constant"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *repository) GetLastPrice(ctx context.Context, symbol string) (float64, error) {
	key := constant.RedisKeyLatestPrice + symbol

	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

func (r *repository) GetRecentTrades(ctx context.Context, symbol string, limit int64) ([]entity.Trade, error) {
	collection := r.mongo.Database("market_db").Collection("trades")

	filter := bson.M{"symbol": symbol}

	opts := options.Find().SetSort(bson.M{"trade_time": -1}).SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	trades := make([]entity.Trade, 0)
	if err = cursor.All(ctx, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *repository) SubscribeTrades(ctx context.Context) <-chan string {
	msgChan := make(chan string)

	go func() {
		pubsub := r.redis.Subscribe(ctx, constant.RedisChannelMarketTrades)
		defer pubsub.Close()

		ch := pubsub.Channel()

		for msg := range ch {
			msgChan <- msg.Payload
		}
		close(msgChan)
	}()

	return msgChan
}
