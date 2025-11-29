package handler

import (
	"context"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/internal/service"
)

type Service interface {
	GetPriceSnapshot(ctx context.Context, symbol string) (*service.PriceSnapshot, error)
	GetTrades(ctx context.Context, symbol string) ([]entity.Trade, error)
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type handler struct {
	service Service
}
