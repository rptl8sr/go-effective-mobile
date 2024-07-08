package app

import (
	"context"
	"fmt"
	"go-effective-mobile/internal/config"
	"go-effective-mobile/internal/logger"
	"go-effective-mobile/internal/router"
	"go-effective-mobile/internal/storage/db"
	"net/http"
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

	err = db.Init(ctx, cfg.DSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &App{
		Port: cfg.Port(),
		Ctx:  ctx,
	}, nil
}

func (a *App) Run() error {
	logger.Info(fmt.Sprintf("Starting server on port %d", a.Port))
	r := router.New()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: r,
	}

	return srv.ListenAndServe()
}
