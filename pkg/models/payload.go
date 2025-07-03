package models

import "fmt"

type WorkMode string

const (
	WorkModeLoadFirst    = "load_first"
	WorkModeBatteryFirst = "battery_first"
	Online               = "online"
	Offline              = "offline"
	Heating              = "heating"
	SmartSelfUse         = "smart_self_use"
	Fault                = "fault"
	OnGrid               = "on_grid"
	OffGrid              = "off_grid"
)

func StatusFromString(s string) string {
	switch s {
	case "-1":
		return Offline
	case "0":
		return WorkModeLoadFirst
	case "1":
		return WorkModeBatteryFirst
	case "2":
		return SmartSelfUse
	case "4":
		return Fault
	case "5":
		return Heating
	case "6":
		return OnGrid
	case "7":
		return OffGrid
	}
	return fmt.Sprintf("invalid_%s", s)
}

func WorkModeFromString(s string) WorkMode {
	if s == "0" {
		return WorkModeLoadFirst
	}
	return WorkModeBatteryFirst
}

func IntFromWorkMode(s WorkMode) int {
	if s == WorkModeLoadFirst {
		return 0
	}
	if s == WorkModeBatteryFirst {
		return 1
	}
	return -1
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
	ChargingLimit        *float64  `json:"charging_limit,omitempty"`
	DischargeLimit       *float64  `json:"discharge_limit,omitempty"`
	DefaultACCouplePower *float64  `json:"default_output_w,omitempty"`
	DefaultMode          *WorkMode `json:"default_mode,omitempty"`
}

func (p *ParameterPayload) UpdateFrom(src ParameterPayload) {
	if src.ChargingLimit != nil {
		p.ChargingLimit = src.ChargingLimit
	}
	if src.DischargeLimit != nil {
		p.DischargeLimit = src.DischargeLimit
	}
	if src.DefaultACCouplePower != nil {
		p.DefaultACCouplePower = src.DefaultACCouplePower
	}
	if src.DefaultMode != nil {
		p.DefaultMode = src.DefaultMode
	}
}

func EmptyParameterPayload() ParameterPayload {
	chargingLimit := 100.0
	dischargeLimit := 10.0
	defaultACCouplePower := 150.0
	var defaultMode WorkMode = WorkModeLoadFirst

	return ParameterPayload{
		ChargingLimit:        &chargingLimit,
		DischargeLimit:       &dischargeLimit,
		DefaultACCouplePower: &defaultACCouplePower,
		DefaultMode:          &defaultMode,
	}
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
