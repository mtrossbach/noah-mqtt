package polling

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/pkg/models"
)

func (s *Service) parametersSubscription(sn string) func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		changedSomething := false
		var payload models.ParameterPayload
		if err := json.Unmarshal(message.Payload(), &payload); err != nil {
			slog.Error("unable to unmarshal parameter command payload", slog.String("error", err.Error()))
		}

		if payload.OutputPower != nil {
			slog.Info("Trying to set default power", slog.String("device", sn), slog.Int("power", int(*payload.OutputPower)))
			if err := s.options.GrowattClient.SetDefaultPower(sn, *payload.OutputPower); err != nil {
				slog.Error("unable to set default power", slog.String("error", err.Error()), slog.String("device", sn))
			} else {
				changedSomething = true
				slog.Info("set default power", slog.String("device", sn), slog.Int("power", int(*payload.OutputPower)))
			}
		}

		if payload.ChargingLimit != nil && payload.DischargeLimit != nil {
			slog.Info("Trying to set charging/discharge limit", slog.String("device", sn), slog.Float64("chargingLimit", *payload.ChargingLimit), slog.Float64("dischargeLimit", *payload.DischargeLimit))
			if err := s.options.GrowattClient.SetSocLimit(sn, *payload.ChargingLimit, *payload.DischargeLimit); err != nil {
				slog.Error("unable to set charging/discharge limit", slog.String("error", err.Error()))
			} else {
				changedSomething = true
				slog.Info("set charging/discharge limit", slog.String("device", sn), slog.Float64("chargingLimit", *payload.ChargingLimit), slog.Float64("dischargeLimit", *payload.DischargeLimit))
			}
		}

		if changedSomething {
			s.pollParameterData(sn)
		}
	}
}
