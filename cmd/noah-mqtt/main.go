package main

import (
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/internal/logging"
	"noah-mqtt/internal/service"
	"os"
	"os/signal"
	"os/user"
	"syscall"
)

func main() {
	cfg := config.Get()
	logging.Init(cfg.LogLevel)
	if err := config.Validate(); err != nil {
		panic(err)
	}

	slog.Info("noah-mqtt started", slog.String("version", cfg.Version))

	if currentUser, err := user.Current(); err == nil {
		slog.Info("running as", slog.String("username", currentUser.Username), slog.String("uid", currentUser.Uid))
	}

	service.Start()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}
