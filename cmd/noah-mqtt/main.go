package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/internal/endpoint_mqtt"
	"noah-mqtt/internal/growatt_app"
	"noah-mqtt/internal/growatt_web"
	"noah-mqtt/internal/homeassistant"
	"noah-mqtt/internal/logging"
	"noah-mqtt/internal/misc"
	"os"
	"os/signal"
	"os/user"
	"strings"
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
		runApp(cfg, client)
	})

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}

func runApp(cfg config.Config, client mqtt.Client) {
	haService := homeassistant.NewService(homeassistant.Options{
		MqttClient:  client,
		TopicPrefix: cfg.HomeAssistant.TopicPrefix,
		Version:     version,
	})

	mqttEndpoint := endpoint_mqtt.NewEndpoint(endpoint_mqtt.Options{
		MqttClient:  client,
		TopicPrefix: cfg.Mqtt.TopicPrefix,
		HaClient:    haService,
	})

	mode := strings.ToLower(strings.TrimSpace(cfg.Growatt.APIMode))
	switch mode {
	case "app":
		slog.Info("setting mode", slog.String("mode", mode))
		growattApp := growatt_app.NewGrowattAppService(growatt_app.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlApp,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
		})

		growattApp.AddEndpoint(mqttEndpoint)
		if err := growattApp.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}
		growattApp.StartPolling()
	case "web":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:       cfg.Growatt.ServerUrlWeb,
			Username:        cfg.Growatt.Username,
			Password:        cfg.Growatt.Password,
			PollingInterval: cfg.PollingInterval,
		})

		growattService.AddEndpoint(mqttEndpoint)
		if err := growattService.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}

		slog.Warn("web mode does not support setting parameters")
		growattService.StartPolling()

	case "web+app":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:       cfg.Growatt.ServerUrlWeb,
			Username:        cfg.Growatt.Username,
			Password:        cfg.Growatt.Password,
			PollingInterval: cfg.PollingInterval,
		})

		growattService.AddEndpoint(mqttEndpoint)
		if err := growattService.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}

		growattApp := growatt_app.NewGrowattAppService(growatt_app.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlApp,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
		})

		mqttEndpoint.SetParameterApplier(growattApp)
		growattService.StartPolling()
	default:
		misc.Panic(fmt.Errorf("invalid growatt api type: %s", cfg.Growatt.APIMode))
	}
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
