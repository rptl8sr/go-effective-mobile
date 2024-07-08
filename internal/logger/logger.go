package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

func Init() {
	// TODO: add options with logging level
	once.Do(func() {
		globalLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
}

func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}
