package growatt_app

import (
	"log/slog"
	"noah-mqtt/pkg/models"
)

func (g *GrowattAppService) pollStatus(device models.NoahDevicePayload) {
	if data, err := g.client.GetNoahStatus(device.Serial); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		payload := devicePayload(data)
		for _, e := range g.endpoints {
			e.PublishDeviceStatus(device, payload)
		}
	}
}

func (g *GrowattAppService) pollBatteryDetails(device models.NoahDevicePayload) {
	if data, err := g.client.GetBatteryData(device.Serial); err != nil {
		slog.Error("could not get battery data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		var batteryPayloads []models.BatteryPayload

		for _, bat := range data.Obj.Batter {
			batteryPayloads = append(batteryPayloads, batteryPayload(&bat))
		}

		for _, e := range g.endpoints {
			e.PublishBatteryDetails(device, batteryPayloads)
		}
	}
}

func (g *GrowattAppService) pollParameterData(device models.NoahDevicePayload) {
	if data, err := g.client.GetNoahInfo(device.Serial); err != nil {
		slog.Error("could not get parameter data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		payload := parameterPayload(data)
		for _, e := range g.endpoints {
			e.PublishParameterData(device, payload)
		}
	}
}
