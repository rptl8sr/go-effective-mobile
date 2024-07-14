package main

import (
	"context"
	"os"

	"go-effective-mobile/internal/app"
	"go-effective-mobile/internal/logger"
)

func main() {
	// TODO: add logging level by using flags
	logger.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.New(ctx)
	if err != nil {
		logger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
