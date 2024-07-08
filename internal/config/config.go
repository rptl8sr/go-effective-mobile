package config

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"go-effective-mobile/internal/logger"
)

type cfg struct {
	port uint16
	db   db
}

type db struct {
	host     string `env:"PG_HOST"`
	port     string `env:"PG_PORT"`
	user     string `env:"PG_USER"`
	password string `env:"PG_PASSWORD"`
	name     string `env:"PG_DATABASE"`
	sslMode  string `env:"PG_SSL_MODE"`
}

type Config interface {
	Port() uint16
	DSN() string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file")
		return nil, err
	}

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		logger.Error("Error parsing PORT")
		return nil, err
	}

	if port < 1 || port > math.MaxUint16 {
		logger.Error("Invalid port number")
		return nil, err
	}

	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DATABASE")
	dbSSLMode := os.Getenv("PG_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		return nil, err
	}

	dbConfig := db{
		host:     dbHost,
		port:     dbPort,
		user:     dbUser,
		password: dbPassword,
		name:     dbName,
		sslMode:  dbSSLMode,
	}

	return &cfg{
		port: uint16(port),
		db:   dbConfig,
	}, nil
}

func (c *cfg) Port() uint16 {
	return c.port
}

func (c *cfg) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.db.host, c.db.port, c.db.name, c.db.user, c.db.password, c.db.sslMode)

	return dsn
}
