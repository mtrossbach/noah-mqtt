package polling

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

func (s *Service) pollStatus(serialNumber string) {
	if data, err := s.options.GrowattClient.GetNoahStatus(serialNumber); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", serialNumber))
	} else {
		if b, err := json.Marshal(devicePayload(data)); err != nil {
			slog.Error("could not marshal device data", slog.String("error", err.Error()), slog.String("device", serialNumber))
		} else {
			s.options.MqttClient.Publish(s.deviceStateTopic(serialNumber), 0, false, string(b))
			slog.Debug("device data received", slog.String("data", string(b)), slog.String("device", serialNumber))
		}
	}
}

func (s *Service) pollBatteryDetails(serialNumber string) {
	if data, err := s.options.GrowattClient.GetBatteryData(serialNumber); err != nil {
		slog.Error("could not get battery data", slog.String("error", err.Error()), slog.String("device", serialNumber))
	} else {
		var logData []any
		for i, bat := range data.Obj.Batter {
			if b, err := json.Marshal(batteryPayload(&bat)); err != nil {
				slog.Error("could not marshal battery data", slog.String("error", err.Error()))
			} else {
				s.options.MqttClient.Publish(s.stateTopicBattery(serialNumber, i), 0, false, string(b))
				logData = append(logData, slog.String(fmt.Sprintf("BAT%d", i), string(b)))
			}
		}
		logData = append(logData, slog.String("device", serialNumber))
		slog.Debug("battery data received", logData...)
	}
}

func (s *Service) pollParameterData(serialNumber string) {
	if data, err := s.options.GrowattClient.GetNoahInfo(serialNumber); err != nil {
		slog.Error("could not get parameter data", slog.String("error", err.Error()), slog.String("device", serialNumber))
	} else {
		if b, err := json.Marshal(parameterPayload(data)); err != nil {
			slog.Error("could not marshal parameter data", slog.String("error", err.Error()), slog.String("device", serialNumber))
		} else {
			s.options.MqttClient.Publish(s.parameterStateTopic(serialNumber), 0, false, string(b))
			slog.Debug("parameter data received", slog.String("data", string(b)), slog.String("device", serialNumber))
		}
	}
}
