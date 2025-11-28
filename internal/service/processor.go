package service

import (
	"context"
	"encoding/json"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/pkg/logger"
)

func (u *service) ProcessTradeEvent(ctx context.Context, payload string) error {
	var trade entity.Trade
	if err := json.Unmarshal([]byte(payload), &trade); err != nil {

		return err
	}

	logger.Log.Infof(ctx, "Processed: %s Price: %.2f", trade.Symbol, trade.Price)

	if err := u.repository.InsertTrade(ctx, trade); err != nil {
		logger.Log.Errorf(ctx, "Failed to insert trade history: %v", err)

	}

	if err := u.repository.UpdateLastPrice(ctx, trade); err != nil {
		logger.Log.Errorf(ctx, "Failed to update cache: %v", err)
	}

	return nil
}
