package handler

import (
	"context"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
	"github.com/ArdyJunata/go-realtime-market-data/internal/service"
)

type Service interface {
	GetPriceSnapshot(ctx context.Context, symbol string) (*service.PriceSnapshot, error)
	GetTrades(ctx context.Context, symbol string) ([]entity.Trade, error)
	GetTradeStream(ctx context.Context) <-chan string
}

func NewHandler(service Service, hub *Hub) *handler {
	return &handler{
		service: service,
		hub:     hub,
	}
}

type handler struct {
	service Service
	hub     *Hub
}
