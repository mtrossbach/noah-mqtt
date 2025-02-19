package growatt_app

import (
	"noah-mqtt/internal/misc"
	"noah-mqtt/pkg/models"
)

func devicePayload(n *NoahStatus) models.DevicePayload {
	return models.DevicePayload{
		OutputPower:           misc.ParseFloat(n.Obj.Pac),
		SolarPower:            misc.ParseFloat(n.Obj.Ppv),
		Soc:                   misc.ParseFloat(n.Obj.Soc),
		ChargePower:           misc.ParseFloat(n.Obj.ChargePower),
		DischargePower:        misc.ParseFloat(n.Obj.DisChargePower),
		BatteryNum:            int(misc.ParseFloat(n.Obj.BatteryNum)),
		GenerationTotalEnergy: misc.ParseFloat(n.Obj.EacTotal),
		GenerationTodayEnergy: misc.ParseFloat(n.Obj.EacToday),
		WorkMode:              models.WorkModeFromString(n.Obj.WorkMode),
		Status:                models.StatusFromString(n.Obj.Status),
	}
}

func batteryPayload(n *BatteryDetails) models.BatteryPayload {
	return models.BatteryPayload{
		SerialNumber: n.SerialNum,
		Soc:          misc.ParseFloat(n.Soc),
		Temperature:  misc.ParseFloat(n.Temp),
	}
}

func parameterPayload(n *NoahInfo) models.ParameterPayload {
	chargingLimit := misc.ParseFloat(n.Obj.Noah.ChargingSocHighLimit)
	dischargeLimit := misc.ParseFloat(n.Obj.Noah.ChargingSocLowLimit)
	outputPower := misc.ParseFloat(n.Obj.Noah.DefaultPower)

	return models.ParameterPayload{
		ChargingLimit:  &chargingLimit,
		DischargeLimit: &dischargeLimit,
		OutputPower:    &outputPower,
	}
}
