package main

import (
	"context"
	"go-effective-mobile/internal/app"
	"go-effective-mobile/internal/logger"
	"os"
)

func main() {
	logger.Init()

	ctx, _ := context.WithCancel(context.Background())

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
