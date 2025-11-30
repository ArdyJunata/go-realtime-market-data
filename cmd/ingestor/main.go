package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/constant"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/database"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/logger"
	"github.com/gorilla/websocket"
)

func main() {
	config.LoadConfig(".env")

	logger.InitLogger()
	ctx := context.Background()

	redis, err := database.InitRedis()
	if err != nil {
		logger.Log.Errorf(ctx, "Redis error: %v", err)
		panic(err)
	}

	for {
		logger.Log.Infof(ctx, "Connecting to Binance: %s", constant.BinanceWs)
		conn, _, err := websocket.DefaultDialer.Dial(constant.BinanceWs, nil)
		if err != nil {
			logger.Log.Errorf(ctx, "Dial error: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Log.Infof(ctx, "Connected to Binance!")

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Errorf(ctx, "Read error: %v", err)
				conn.Close()
				break
			}

			var raw entity.BinanceAggTrade
			if err := json.Unmarshal(message, &raw); err != nil {
				logger.Log.Warnf(ctx, "Bad JSON: %v", err)
				continue
			}

			price, _ := strconv.ParseFloat(raw.Price, 64)
			qty, _ := strconv.ParseFloat(raw.Quantity, 64)

			cleanTrade := entity.Trade{
				ID:        strconv.FormatInt(raw.TradeID, 10),
				Symbol:    raw.Symbol,
				Price:     price,
				Quantity:  qty,
				TradeTime: time.UnixMilli(raw.TradeTime),
				IsBuyer:   !raw.IsBuyer,
			}

			payload, _ := json.Marshal(cleanTrade)
			if err := redis.Publish(ctx, constant.RedisChannelMarketTrades, payload).Err(); err != nil {
				logger.Log.Errorf(ctx, "Redis Publish Failed: %v", err)
			}

			logger.Log.Infof(ctx, "Tick: %.2f", cleanTrade.Price)
		}

		time.Sleep(1 * time.Second)
	}
}
