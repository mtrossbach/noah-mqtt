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
	Version     string
}

type Service struct {
	options Options

	devices           []DeviceInfo
	statusChangeToken mqtt.Token
}

func NewService(opts Options) *Service {
	s := &Service{
		options: opts,
	}
	s.statusChangeToken = opts.MqttClient.Subscribe(fmt.Sprintf("%s/status", opts.TopicPrefix), 0, s.haStatusChange)

	return s
}

func (s *Service) haStatusChange(client mqtt.Client, message mqtt.Message) {
	s.sendDiscovery()
}

func (s *Service) SetDevices(devices []DeviceInfo) {
	s.devices = devices
	s.sendDiscovery()
}

func (s *Service) sendDiscovery() {
	for _, d := range s.devices {
		sensors := generateSensorDiscoveryPayload(s.options.Version, d)
		for _, sensor := range sensors {
			if b, err := json.Marshal(sensor); err != nil {
				slog.Error("could not marshal sensor discovery payload", slog.Any("sensor", sensor))
			} else {
				topic := s.sensorTopic(sensor)
				s.options.MqttClient.Publish(topic, 0, false, string(b))
			}
		}

		numbers := generateNumberDiscoveryPayload(s.options.Version, d)
		for _, number := range numbers {
			if b, err := json.Marshal(number); err != nil {
				slog.Error("could not marshal number discovery payload", slog.Any("number", number))
			} else {
				topic := s.numberTopic(number)
				s.options.MqttClient.Publish(topic, 0, false, string(b))
			}
		}
	}
}

func (s *Service) sensorTopic(sensor Sensor) string {
	return fmt.Sprintf("%s/sensor/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("noah_%s", sensor.Device.SerialNumber), strings.ReplaceAll(sensor.Name, " ", ""))
}

func (s *Service) numberTopic(number Number) string {
	return fmt.Sprintf("%s/number/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("noah_%s", number.Device.SerialNumber), strings.ReplaceAll(number.Name, " ", ""))

}
