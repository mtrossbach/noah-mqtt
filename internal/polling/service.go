package polling

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/internal/growatt"
	"noah-mqtt/internal/homeassistant"
	"os"
	"time"
)

type Options struct {
	GrowattClient     *growatt.Client
	HaClient          *homeassistant.Service
	MqttClient        mqtt.Client
	PollingInterval   time.Duration
	TopicPrefix       string
	DetailsCycleSkips int
}

type Service struct {
	options       Options
	serialNumbers []string
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

	s.enumerateDevices()

	go s.poll()
}

func (s *Service) enumerateDevices() {
	serialNumbers := s.fetchNoahSerialNumbers()
	var devices []homeassistant.DeviceInfo
	for _, serialNumber := range serialNumbers {
		if data, err := s.options.GrowattClient.GetNoahInfo(serialNumber); err != nil {
			slog.Error("could not noah status", slog.String("error", err.Error()), slog.String("serialNumber", serialNumber))
		} else {
			batCount := len(data.Obj.Noah.BatSns)
			var batteries []homeassistant.BatteryInfo
			for i := 0; i < batCount; i++ {
				batteries = append(batteries, homeassistant.BatteryInfo{
					Alias:      fmt.Sprintf("BAT%d", i),
					StateTopic: s.stateTopicBattery(serialNumber, i),
				})
			}

			devices = append(devices, homeassistant.DeviceInfo{
				SerialNumber:          serialNumber,
				Model:                 data.Obj.Noah.Model,
				Version:               data.Obj.Noah.Version,
				Alias:                 data.Obj.Noah.Alias,
				StateTopic:            s.deviceStateTopic(serialNumber),
				ParameterStateTopic:   s.parameterStateTopic(serialNumber),
				ParameterCommandTopic: s.parameterCommandTopic(serialNumber),
				Batteries:             batteries,
			})
		}
	}

	for _, sn := range s.serialNumbers {
		s.options.MqttClient.Unsubscribe(s.parameterCommandTopic(sn))
	}

	s.serialNumbers = serialNumbers

	for _, sn := range s.serialNumbers {
		s.options.MqttClient.Subscribe(s.parameterCommandTopic(sn), 0, s.parametersSubscription(sn))
	}

	s.options.HaClient.SetDevices(devices)
}

func (s *Service) deviceStateTopic(serialNumber string) string {
	return fmt.Sprintf("%s/%s", s.options.TopicPrefix, serialNumber)
}

func (s *Service) stateTopicBattery(serialNumber string, index int) string {
	return fmt.Sprintf("%s/%s/BAT%d", s.options.TopicPrefix, serialNumber, index)
}

func (s *Service) parameterStateTopic(serialNumber string) string {
	return fmt.Sprintf("%s/%s/parameters", s.options.TopicPrefix, serialNumber)
}

func (s *Service) parameterCommandTopic(serialNumber string) string {
	return fmt.Sprintf("%s/%s/parameters/set", s.options.TopicPrefix, serialNumber)
}

func (s *Service) fetchNoahSerialNumbers() []string {
	slog.Info("fetching plant list")
	list, err := s.options.GrowattClient.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		panic(err)
	}

	var serialNumbers []string

	for _, plant := range list.PlantList {
		slog.Info("fetch plant details", slog.Int("plantId", plant.ID))
		if info, err := s.options.GrowattClient.GetNoahPlantInfo(fmt.Sprintf("%d", plant.ID)); err != nil {
			slog.Error("could not get plant info", slog.Int("plantId", plant.ID), slog.String("error", err.Error()))
		} else {
			if len(info.Obj.DeviceSn) > 0 {
				serialNumbers = append(serialNumbers, info.Obj.DeviceSn)
				slog.Info("found device sn", slog.String("deviceSn", info.Obj.DeviceSn), slog.Int("plantId", plant.ID), slog.String("topic", s.deviceStateTopic(info.Obj.DeviceSn)))
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

func (s *Service) poll() {
	slog.Info("start polling growatt", slog.Int("interval", int(s.options.PollingInterval/time.Second)))
	i := 0
	for {
		for _, serialNumber := range s.serialNumbers {
			s.pollStatus(serialNumber)
			if i%(s.options.DetailsCycleSkips+1) == 0 {
				s.pollBatteryDetails(serialNumber)
				s.pollParameterData(serialNumber)
			}
		}
		<-time.After(s.options.PollingInterval)

		i += 1
		if i >= 1000 {
			i = 0
		}
	}
}
