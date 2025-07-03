package homeassistant

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type HaClient interface {
	SetDevices(devices []DeviceInfo)
}

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
	go s.discoveryLooper()
	return s
}

func (s *Service) discoveryLooper() {
	for {
		<-time.After(6 * time.Hour)
		if len(s.devices) > 0 {
			s.sendDiscovery()
		}
	}
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

		selects := generateSelectDiscoveryPayload(s.options.Version, d)
		for _, sel := range selects {
			if b, err := json.Marshal(sel); err != nil {
				slog.Error("could not marshal select discovery payload", slog.Any("select", sel))
			} else {
				topic := s.selectTopic(sel)
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

		binarySensors := generateBinarySensorDiscoveryPayload(s.options.Version, d)
		for _, sensor := range binarySensors {
			if b, err := json.Marshal(sensor); err != nil {
				slog.Error("could not marshal binary sensor discovery payload", slog.Any("sensor", sensor))
			} else {
				topic := s.binarySensorTopic(sensor)
				s.options.MqttClient.Publish(topic, 0, false, string(b))
			}
		}
	}
}

func (s *Service) sensorTopic(sensor Sensor) string {
	return fmt.Sprintf("%s/sensor/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("nexa_%s", sensor.Device.SerialNumber), strings.ReplaceAll(sensor.Name, " ", ""))
}

func (s *Service) selectTopic(sensor Select) string {
	return fmt.Sprintf("%s/select/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("nexa_%s", sensor.Device.SerialNumber), strings.ReplaceAll(sensor.Name, " ", ""))
}

func (s *Service) binarySensorTopic(sensor BinarySensor) string {
	return fmt.Sprintf("%s/binary_sensor/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("nexa_%s", sensor.Device.SerialNumber), strings.ReplaceAll(sensor.Name, " ", ""))
}

func (s *Service) numberTopic(number Number) string {
	return fmt.Sprintf("%s/number/%s/%s/config", s.options.TopicPrefix, fmt.Sprintf("nexa_%s", number.Device.SerialNumber), strings.ReplaceAll(number.Name, " ", ""))

}
