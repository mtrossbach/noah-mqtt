package polling

import (
	"noah-mqtt/internal/growatt"
	"noah-mqtt/pkg/models"
	"strconv"
)

func noahStatusToPayload(n *growatt.NoahStatus) models.Payload {
	return models.Payload{
		OutputPower:           parseFloat(n.Obj.Pac),
		SolarPower:            parseFloat(n.Obj.Ppv),
		Soc:                   parseFloat(n.Obj.Soc),
		ChargePower:           parseFloat(n.Obj.ChargePower),
		DischargePower:        parseFloat(n.Obj.DisChargePower),
		BatteryNum:            int(parseFloat(n.Obj.BatteryNum)),
		GenerationTotalEnergy: parseFloat(n.Obj.EacTotal),
		GenerationTodayEnergy: parseFloat(n.Obj.EacToday),
		WorkMode:              workModeFromString(n.Obj.WorkMode),
	}
}

func workModeFromString(s string) models.WorkMode {
	if s == "0" {
		return models.WorkModeLoadFirst
	}
	return models.WorkModeBatteryFirst
}

func parseFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	} else {
		return 0
	}
}
