package service

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/config"
	"noah-mqtt/pkg/models"
	"strings"
)

func SendDiscovery(mqttClient mqtt.Client, cfg config.HomeAssistant, sensorTopic string, name string, serial string) {
	dev := models.HomeAssistantDevice{
		Identifiers:  []string{fmt.Sprintf("noah_%s", serial)},
		Name:         name,
		Manufacturer: "Growatt",
		Model:        "Noah",
		SerialNumber: serial,
	}

	topic := fmt.Sprintf("%s/sensor/noah_%s", cfg.TopicPrefix, serial)

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Output Power",
		DeviceClass:       "power",
		StateClass:        "measurement",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "W",
		ValueTemplate:     "{{ value_json.output_w }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "output_power"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Solar Power",
		Icon:              "mdi:solar-power",
		DeviceClass:       "power",
		StateClass:        "measurement",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "W",
		ValueTemplate:     "{{ value_json.solar_w }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "solar_power"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Charging Power",
		Icon:              "mdi:battery-plus",
		DeviceClass:       "power",
		StateClass:        "measurement",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "W",
		ValueTemplate:     "{{ value_json.charge_w }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "charging_power"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Discharge Power",
		Icon:              "mdi:battery-minus",
		DeviceClass:       "power",
		StateClass:        "measurement",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "W",
		ValueTemplate:     "{{ value_json.discharge_w }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "discharge_power"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Generation Total",
		DeviceClass:       "energy",
		StateClass:        "total_increasing",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "kWh",
		ValueTemplate:     "{{ value_json.generation_total_kwh }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "generation_total"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "Generation Today",
		DeviceClass:       "energy",
		StateClass:        "total_increasing",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "kWh",
		ValueTemplate:     "{{ value_json.generation_today_kwh }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "generation_today"),
		Device:            dev,
	})

	sendDiscoveryMessage(mqttClient, topic, models.HomeAssistantSensor{
		Name:              "SoC",
		DeviceClass:       "battery",
		StateClass:        "measurement",
		StateTopic:        sensorTopic,
		UnitOfMeasurement: "%",
		ValueTemplate:     "{{ value_json.soc }}",
		UniqueId:          fmt.Sprintf("%s_%s", serial, "soc"),
		Device:            dev,
	})

}

func sendDiscoveryMessage(mqttClient mqtt.Client, baseTopic string, sensor models.HomeAssistantSensor) {
	topic := fmt.Sprintf("%s/%s/config", baseTopic, strings.ReplaceAll(sensor.Name, " ", ""))
	if b, err := json.Marshal(sensor); err != nil {
		slog.Error("could not marshal sensor discovery payload", slog.Any("sensor", sensor))
	} else {
		mqttClient.Publish(topic, 1, true, string(b))
	}

}
