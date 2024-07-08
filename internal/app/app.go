package app

import (
	"context"
	"go-effective-mobile/internal/config"
)

type App struct {
	Port uint16
	Ctx  context.Context
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &App{
		Port: cfg.Port(),
		Ctx:  ctx,
	}, nil
}
