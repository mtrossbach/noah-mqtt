package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/internal/growatt"
	"noah-mqtt/internal/homeassistant"
	"noah-mqtt/internal/logging"
	"noah-mqtt/internal/misc"
	"noah-mqtt/internal/polling"
	"os"
	"os/signal"
	"os/user"
	"syscall"
)

var (
	version = "local"
	commit  = "none"
)

func main() {
	cfg := config.Get()
	logging.Init(cfg.LogLevel)
	if err := config.Validate(); err != nil {
		slog.Error("couldn't validate config", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	slog.Info("noah-mqtt started", slog.String("version", version), slog.String("commit", commit))

	if currentUser, err := user.Current(); err == nil {
		slog.Info("running as", slog.String("username", currentUser.Username), slog.String("uid", currentUser.Uid))
	}

	connectMqtt(cfg.Mqtt, func(client mqtt.Client) {
		growattClient := growatt.NewClient(cfg.Growatt.ServerUrl, cfg.Growatt.Username, cfg.Growatt.Password)
		haService := homeassistant.NewService(homeassistant.Options{
			MqttClient:  client,
			TopicPrefix: cfg.HomeAssistant.TopicPrefix,
			Version:     version,
		})
		pollingService := polling.NewService(polling.Options{
			GrowattClient:                 growattClient,
			HaClient:                      haService,
			MqttClient:                    client,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
			TopicPrefix:                   cfg.Mqtt.TopicPrefix,
		})
		pollingService.Start()
	})

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}

func connectMqtt(mqttCfg config.Mqtt, onConnected func(client mqtt.Client)) {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", mqttCfg.Host, mqttCfg.Port)).
		SetClientID(mqttCfg.ClientId).
		SetUsername(mqttCfg.Username).
		SetPassword(mqttCfg.Password)

	opts.OnConnect = func(client mqtt.Client) {
		slog.Info("connected to mqtt broker")
		onConnected(client)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		slog.Error("lost connection to mqtt broker", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	c := mqtt.NewClient(opts)
	slog.Info("connecting to mqtt broker", slog.String("host", mqttCfg.Host), slog.Int("port", mqttCfg.Port), slog.String("clientId", mqttCfg.ClientId), slog.String("username", mqttCfg.Username))
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("could not connect to mqtt broker", slog.String("error", token.Error().Error()))
		misc.Panic(token.Error())
	}
}
