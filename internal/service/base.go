package service

import (
	"context"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
)

type Repository interface {
	InsertTrade(ctx context.Context, trade entity.Trade) error
	UpdateLastPrice(ctx context.Context, trade entity.Trade) error
	GetLastPrice(ctx context.Context, symbol string) (float64, error)
	GetRecentTrades(ctx context.Context, symbol string, limit int64) ([]entity.Trade, error)
	SubscribeTrades(ctx context.Context) <-chan string
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}

type PriceSnapshot struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
