package main

import (
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	slog.Info("noah-mqtt started", slog.String("version", cfg.Version))
	service.Start()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}
