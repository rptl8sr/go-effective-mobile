package config

import (
	"github.com/joho/godotenv"
	"go-effective-mobile/internal/logger"
	"math"
	"os"
	"strconv"
)

const (
	defPort = 9090
)

type cfg struct {
	port uint16
}

type Config interface {
	Port() uint16
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", "error", err.Error())
		return nil, err
	}

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		logger.Error("Error parsing PORT", "error", err.Error())
		return nil, err
	}

	if port < 1 || port > math.MaxUint16 {
		logger.Error("Invalid port number")
		return nil, err
	}

	return &cfg{
		port: uint16(port),
	}, nil
}

func (c *cfg) Port() uint16 {
	return c.port
}
