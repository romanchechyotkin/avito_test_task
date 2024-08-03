package logger

import (
	"log/slog"
	"os"
	"strings"
)

const defaultLevel = slog.LevelDebug

// todo refactor => devLogger() prodLogger()
// todo add source for dev
func New() *slog.Logger {
	var handler slog.Handler

	if env := os.Getenv("APP_ENV"); env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: configLevel(),
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: configLevel(),
		})
	}

	return slog.New(handler)
}

func configLevel() slog.Level {
	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = defaultLevel
	}

	return logLevel
}
