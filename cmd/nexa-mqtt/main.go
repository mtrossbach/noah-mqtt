package main

import (
	"fmt"
	"log/slog"
	"nexa-mqtt/internal/config"
	"nexa-mqtt/internal/endpoint_mqtt"
	"nexa-mqtt/internal/growatt_app"
	"nexa-mqtt/internal/growatt_web"
	"nexa-mqtt/internal/homeassistant"
	"nexa-mqtt/internal/logging"
	"nexa-mqtt/internal/misc"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	slog.Info("nexa-mqtt started", slog.String("version", version), slog.String("commit", commit))

	if currentUser, err := user.Current(); err == nil {
		slog.Info("running as", slog.String("username", currentUser.Username), slog.String("uid", currentUser.Uid))
	}

	app := NewApp(cfg)
	connectMqtt(cfg.Mqtt, app)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}

type App struct {
	mode           string
	cfg            config.Config
	growattService *growatt_web.GrowattService
	growattApp     *growatt_app.GrowattAppService
}

func (a *App) onMqttDisconnect() {
	if a.growattService != nil {
		a.growattService.StopPolling()
		a.growattService.SetEndpoint(nil)
	}
	if a.growattApp != nil {
		a.growattApp.StopPolling()
		a.growattApp.SetEndpoint(nil)
	}
}

func (a *App) onMqttConnect(client mqtt.Client) {
	haService := homeassistant.NewService(homeassistant.Options{
		MqttClient:  client,
		TopicPrefix: a.cfg.HomeAssistant.TopicPrefix,
		Version:     version,
	})

	mqttEndpoint := endpoint_mqtt.NewEndpoint(endpoint_mqtt.Options{
		MqttClient:  client,
		TopicPrefix: a.cfg.Mqtt.TopicPrefix,
		HaClient:    haService,
	})

	switch a.mode {
	case "app":
		a.growattApp.SetEndpoint(mqttEndpoint)
		a.growattApp.StartPolling()
		mqttEndpoint.SetParameterApplier(a.growattApp)

	case "web":
		a.growattService.SetEndpoint(mqttEndpoint)
		a.growattService.StartPolling()

	case "web+app":
		a.growattService.SetEndpoint(mqttEndpoint)
		a.growattService.StartPolling()
		mqttEndpoint.SetParameterApplier(a.growattApp)
	}
}

func NewApp(cfg config.Config) *App {

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

		if err := growattApp.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}
		return &App{
			mode:       mode,
			cfg:        cfg,
			growattApp: growattApp,
		}

	case "web":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:       cfg.Growatt.ServerUrlWeb,
			Username:        cfg.Growatt.Username,
			Password:        cfg.Growatt.Password,
			PollingInterval: cfg.PollingInterval,
		})

		if err := growattService.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}

		slog.Warn("web mode does not support setting parameters")
		return &App{
			mode:           mode,
			cfg:            cfg,
			growattService: growattService,
		}

	case "web+app":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:       cfg.Growatt.ServerUrlWeb,
			Username:        cfg.Growatt.Username,
			Password:        cfg.Growatt.Password,
			PollingInterval: cfg.PollingInterval,
		})

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

		return &App{
			mode:           mode,
			cfg:            cfg,
			growattService: growattService,
			growattApp:     growattApp,
		}

	default:
		misc.Panic(fmt.Errorf("invalid growatt api type: %s", cfg.Growatt.APIMode))
		return nil
	}
}

func connectMqtt(mqttCfg config.Mqtt, app *App) {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", mqttCfg.Host, mqttCfg.Port)).
		SetClientID(mqttCfg.ClientId).
		SetUsername(mqttCfg.Username).
		SetPassword(mqttCfg.Password)

	opts.OnConnect = func(client mqtt.Client) {
		slog.Info("connected to mqtt broker")
		app.onMqttConnect(client)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		slog.Warn("lost connection to mqtt broker", slog.String("error", err.Error()))
		app.onMqttDisconnect()
	}

	c := mqtt.NewClient(opts)
	slog.Info("connecting to mqtt broker", slog.String("host", mqttCfg.Host), slog.Int("port", mqttCfg.Port), slog.String("clientId", mqttCfg.ClientId), slog.String("username", mqttCfg.Username))
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("could not connect to mqtt broker", slog.String("error", token.Error().Error()))
		misc.Panic(token.Error())
	}
}
