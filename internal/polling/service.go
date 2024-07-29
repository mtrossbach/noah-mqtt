package polling

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/growatt"
	"noah-mqtt/internal/homeassistant"
	"os"
	"time"
)

type Options struct {
	GrowattClient   *growatt.Client
	HaClient        *homeassistant.Service
	MqttClient      mqtt.Client
	PollingInterval time.Duration
	TopicPrefix     string
}

type Service struct {
	options Options
}

func NewService(options Options) *Service {
	return &Service{
		options: options,
	}
}

func (s *Service) Start() {
	if err := s.options.GrowattClient.Login(); err != nil {
		slog.Error("could not login to growatt account", slog.String("error", err.Error()))
		panic(err)
	}

	serialNumbers := s.fetchNoahSerialNumbers()
	for _, serialNumber := range serialNumbers {
		if data, err := s.options.GrowattClient.GetNoahStatus(serialNumber); err != nil {
			slog.Error("could not get data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
		} else {
			sensors := s.options.HaClient.GenerateSensorDiscoveryPayload(data.Obj.Alias, serialNumber, s.stateTopic(serialNumber))
			s.options.HaClient.SendSensorDiscoveryPayload(sensors)
		}
	}

	go s.poll(serialNumbers)
}

func (s *Service) stateTopic(serialNumber string) string {
	return fmt.Sprintf("%s/%s", s.options.TopicPrefix, serialNumber)
}

func (s *Service) fetchNoahSerialNumbers() []string {
	slog.Info("fetching plant list")
	list, err := s.options.GrowattClient.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		panic(err)
	}

	var serialNumbers []string

	for _, plant := range list.Back.Data {
		slog.Info("fetch plant details", slog.String("plantId", plant.PlantID))
		if info, err := s.options.GrowattClient.GetNoahPlantInfo(plant.PlantID); err != nil {
			slog.Error("could not get plant info", slog.String("plantId", plant.PlantID), slog.String("error", err.Error()))
		} else {
			if len(info.Obj.DeviceSn) > 0 {
				serialNumbers = append(serialNumbers, info.Obj.DeviceSn)
				slog.Info("found device sn", slog.String("deviceSn", info.Obj.DeviceSn), slog.String("plantId", plant.PlantID), slog.String("topic", s.stateTopic(info.Obj.DeviceSn)))
			}
		}
	}

	if len(serialNumbers) == 0 {
		slog.Info("no noah devices found")
		<-time.After(60 * time.Second)
		os.Exit(0)
	}

	return serialNumbers
}

func (s *Service) poll(serialNumbers []string) {
	slog.Info("start polling growatt", slog.Int("interval", int(s.options.PollingInterval/time.Second)))
	for {
		for _, serialNumber := range serialNumbers {
			if data, err := s.options.GrowattClient.GetNoahStatus(serialNumber); err != nil {
				slog.Error("could not get data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
			} else {
				if b, err := json.Marshal(noahStatusToPayload(data)); err != nil {
					slog.Error("could not marshal data", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
				} else {
					s.options.MqttClient.Publish(s.stateTopic(serialNumber), 1, true, string(b))
					slog.Debug("data received", slog.String("data", string(b)), slog.String("serialNumber", serialNumber))
				}
			}
		}

		<-time.After(s.options.PollingInterval)
	}
}
