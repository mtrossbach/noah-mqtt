package logging

import (
	"log/slog"
	"os"
	"strings"
)

var (
	logLevel *slog.LevelVar
	logger   *slog.Logger
)

func Init(defaultLogLevel string) {
	logLevel = new(slog.LevelVar)
	jsonHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger = slog.New(jsonHandler)
	slog.SetDefault(logger)

	if err := setLogLevel(defaultLogLevel); err != nil {
		slog.Error("failed to set default log level", slog.String("error", err.Error()))
	}
}

func setLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		slog.Error("invalid log level", slog.String("level", level))
	}
	slog.Info("set log level", slog.String("level", level))
	return nil
}

func IsDebug() bool {
	return logLevel.Level() == slog.LevelDebug
}
