package service

import (
	"context"

	"github.com/ArdyJunata/go-realtime-market-data/internal/entity"
)

type Repository interface {
	InsertTrade(ctx context.Context, trade entity.Trade) error
	UpdateLastPrice(ctx context.Context, trade entity.Trade) error
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}
