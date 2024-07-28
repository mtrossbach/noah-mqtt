package service

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/internal/growatt"
	"time"
)

func Start() {
	connectMqtt(config.Get().Mqtt)
}

func fetchSerialNumbers(client *growatt.Client) []string {
	list, err := client.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		panic(err)
	}

	var serialNumbers []string

	for _, plant := range list.Back.Data {
		if info, err := client.GetNoahPlantInfo(plant.PlantID); err != nil {
			slog.Error("could not get plant info", slog.String("plantId", plant.PlantID), slog.String("error", err.Error()))
		} else {
			if len(info.Obj.DeviceSn) > 0 {
				serialNumbers = append(serialNumbers, info.Obj.DeviceSn)
				slog.Info("found device sn", slog.String("deviceSn", info.Obj.DeviceSn), slog.String("plantId", plant.PlantID))
			}
		}
	}

	return serialNumbers
}

func poll(mqttClient mqtt.Client) {
	cfg := config.Get()
	if len(cfg.Growatt.Username) == 0 || len(cfg.Growatt.Password) == 0 {
		panic("growatt username or password is empty")
	}
	client := growatt.NewClient(cfg.Growatt.Username, cfg.Growatt.Password)
	slog.Info("start polling growatt", slog.String("username", cfg.Growatt.Username))
	_ = client.Login()

	serialNumbers := fetchSerialNumbers(client)
	for _, serialNumber := range serialNumbers {
		if data, err := client.GetNoahStatus(serialNumber); err != nil {
			slog.Error("could not get data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
		} else {
			sensorTopic := fmt.Sprintf("%s/%s", cfg.Mqtt.TopicPrefix, serialNumber)
			SendDiscovery(mqttClient, cfg.HomeAssistant, sensorTopic, data.Obj.Alias, serialNumber)
		}
	}

	for {
		for _, serialNumber := range serialNumbers {
			if data, err := client.GetNoahStatus(serialNumber); err != nil {
				slog.Error("could not get data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
			} else {
				if b, err := json.Marshal(data.ToPayload()); err != nil {
					slog.Error("could not marshal data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
				} else {
					mqttClient.Publish(fmt.Sprintf("%s/%s", cfg.Mqtt.TopicPrefix, serialNumber), 1, true, string(b))
					slog.Info("publish data", slog.String("data", string(b)), slog.String("topic", cfg.Mqtt.TopicPrefix), slog.String("serialNumber", serialNumber))
				}
			}
		}

		<-time.After(cfg.PollingInterval)
	}
}

func connectMqtt(mqttCfg config.Mqtt) {
	if len(mqttCfg.Host) == 0 || mqttCfg.Port < 1 {
		panic("mqtt host or port is empty")
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", mqttCfg.Host, mqttCfg.Port)).
		SetClientID(mqttCfg.ClientId).
		SetUsername(mqttCfg.Username).
		SetUsername(mqttCfg.Password)

	opts.OnConnectionLost = connectLostHandler
	opts.OnConnect = connectHandler
	c := mqtt.NewClient(opts)
	slog.Info("connecting to mqtt broker", slog.String("host", mqttCfg.Host), slog.Int("port", mqttCfg.Port), slog.String("clientId", mqttCfg.ClientId), slog.String("username", mqttCfg.Username))
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	slog.Info("connected to mqtt broker")
	poll(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	slog.Error("lost connection to mqtt broker", slog.String("error", err.Error()))
	panic(err)
}
