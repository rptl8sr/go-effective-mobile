package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go-effective-mobile/internal/logger"
)

var (
	client *Client
)

type Client struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Init(ctx context.Context, dsn string) error {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("pgxpool.ParseConfig error")
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error("pgxpool.NewWithConfig error")
		return err
	}

	pgConnString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		cfg.ConnConfig.Host, cfg.ConnConfig.Port, cfg.ConnConfig.User,
		cfg.ConnConfig.Password, cfg.ConnConfig.Database)

	db, err := sql.Open("pgx", pgConnString)
	if err != nil {
		logger.Error("Cannot open database")
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		logger.Error("goose.SetDialect error")
		return err
	}

	if err = goose.Up(db, "migrations"); err != nil {
		logger.Error("goose.Up error", "error", err)
		_ = db.Close()
		return err
	}

	_ = db.Close()

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
