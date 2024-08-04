package models

type WorkMode string

const (
	WorkModeLoadFirst    = "load_first"
	WorkModeBatteryFirst = "battery_first"
)

type DevicePayload struct {
	OutputPower           float64 `json:"output_w"`
	SolarPower            float64 `json:"solar_w"`
	Soc                   float64 `json:"soc"`
	ChargePower           float64 `json:"charge_w"`
	DischargePower        float64 `json:"discharge_w"`
	BatteryNum            int     `json:"battery_num"`
	GenerationTotalEnergy float64 `json:"generation_total_kwh"`
	GenerationTodayEnergy float64 `json:"generation_today_kwh"`
}

type BatteryPayload struct {
	SerialNumber string  `json:"serial"`
	Soc          float64 `json:"soc"`
	Temperature  float64 `json:"temp"`
}

type ParameterPayload struct {
	ChargingLimit  float64  `json:"charging_limit"`
	DischargeLimit float64  `json:"discharge_limit"`
	OutputPower    float64  `json:"output_power_w"`
	WorkMode       WorkMode `json:"work_mode,omitempty"`
}
