package homeassistant

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"strings"
)

type Options struct {
	MqttClient  mqtt.Client
	TopicPrefix string
}

type Service struct {
	options Options
}

func NewService(opts Options) *Service {
	return &Service{
		options: opts,
	}
}

func (s *Service) SendSensorDiscoveryPayload(sensors []Sensor) {
	for _, sensor := range sensors {
		if b, err := json.Marshal(sensor); err != nil {
			slog.Error("could not marshal sensor discovery payload", slog.Any("sensor", sensor))
		} else {
			topic := s.sensorTopic(sensor)
			s.options.MqttClient.Publish(topic, 1, false, string(b))
		}
	}
}

func (s *Service) sensorTopic(sensor Sensor) string {
	return fmt.Sprintf("%s/sensor/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("noah_%s", sensor.Device.SerialNumber), strings.ReplaceAll(sensor.Name, " ", ""))
}
