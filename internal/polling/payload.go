package polling

import (
	"noah-mqtt/internal/growatt"
	"noah-mqtt/pkg/models"
	"strconv"
)

func devicePayload(n *growatt.NoahStatus) models.DevicePayload {
	return models.DevicePayload{
		OutputPower:           parseFloat(n.Obj.Pac),
		SolarPower:            parseFloat(n.Obj.Ppv),
		Soc:                   parseFloat(n.Obj.Soc),
		ChargePower:           parseFloat(n.Obj.ChargePower),
		DischargePower:        parseFloat(n.Obj.DisChargePower),
		BatteryNum:            int(parseFloat(n.Obj.BatteryNum)),
		GenerationTotalEnergy: parseFloat(n.Obj.EacTotal),
		GenerationTodayEnergy: parseFloat(n.Obj.EacToday),
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

func batteryPayload(n *growatt.BatteryDetails) models.BatteryPayload {
	return models.BatteryPayload{
		SerialNumber: n.SerialNum,
		Soc:          parseFloat(n.Soc),
		Temperature:  parseFloat(n.Temp),
	}
}

func parameterPayload(n *growatt.NoahInfo, workMode string) models.ParameterPayload {
	return models.ParameterPayload{
		ChargingLimit:  parseFloat(n.Obj.Noah.ChargingSocHighLimit),
		DischargeLimit: parseFloat(n.Obj.Noah.ChargingSocLowLimit),
		OutputPower:    parseFloat(n.Obj.Noah.DefaultPower),
		WorkMode:       workModeFromString(workMode),
	}
}
