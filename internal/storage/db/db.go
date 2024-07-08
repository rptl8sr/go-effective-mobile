package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"

	"go-effective-mobile/internal/logger"
)

var (
	client *Client
)

type Client struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func New(ctx context.Context, dsn string) error {
	// TODO: create dsn or get it from cfg
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("pgxpool.ParseConfig error")
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error("pgxpool.NewWithConfig error")
		return err
	}

	client = &Client{
		Pool: pool,
		Ctx:  ctx,
	}
	return nil
}

func Close() {
	if client.Pool != nil {
		logger.Debug("close db connection")
		client.Pool.Close()
	}
}

func Ping() error {
	return client.Pool.Ping(client.Ctx)
}
