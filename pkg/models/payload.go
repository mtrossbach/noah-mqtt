package models

type WorkMode string

const (
	WorkModeLoadFirst    = "load_first"
	WorkModeBatteryFirst = "battery_first"
	Online               = "online"
	Offline              = "offline"
	Heating              = "heating"
)

func StatusFromString(s string) string {
	if s == "1" {
		return Online
	}
	if s == "5" {
		return Heating
	}
	return Offline
}

func WorkModeFromString(s string) WorkMode {
	if s == "0" {
		return WorkModeLoadFirst
	}
	return WorkModeBatteryFirst
}

type DevicePayload struct {
	OutputPower           float64  `json:"output_w"`
	SolarPower            float64  `json:"solar_w"`
	Soc                   float64  `json:"soc"`
	ChargePower           float64  `json:"charge_w"`
	DischargePower        float64  `json:"discharge_w"`
	BatteryNum            int      `json:"battery_num"`
	GenerationTotalEnergy float64  `json:"generation_total_kwh"`
	GenerationTodayEnergy float64  `json:"generation_today_kwh"`
	WorkMode              WorkMode `json:"work_mode,omitempty"`
	Status                string   `json:"status,omitempty"`
}

type BatteryPayload struct {
	SerialNumber string  `json:"serial"`
	Soc          float64 `json:"soc"`
	Temperature  float64 `json:"temp"`
}

type ParameterPayload struct {
	ChargingLimit  *float64 `json:"charging_limit,omitempty"`
	DischargeLimit *float64 `json:"discharge_limit,omitempty"`
	OutputPower    *float64 `json:"output_power_w,omitempty"`
}

type NoahDevicePayload struct {
	PlantId   int                        `json:"plant_id"`
	Serial    string                     `json:"serial"`
	Model     string                     `json:"model"`
	Version   string                     `json:"version"`
	Alias     string                     `json:"alias"`
	Batteries []NoahDeviceBatteryPayload `json:"batteries"`
}

type NoahDeviceBatteryPayload struct {
	Alias string `json:"alias"`
}
