package homeassistant

import (
	"fmt"
	"nexa-mqtt/pkg/models"
)

func generateSensorDiscoveryPayload(appVersion string, info DeviceInfo) []Sensor {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	sensors := []Sensor{
		{
			CommonConfig: CommonConfig{
				Name:        "Output Power",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "output_power"),
				DeviceClass: DeviceClassPower,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.output_w }}",
			},
			StateClass:        StateClassMeasurement,
			UnitOfMeasurement: UnitWatt,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Solar Power",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "solar_power"),
				Icon:        IconSolarPower,
				DeviceClass: DeviceClassPower,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.solar_w }}",
			},
			StateClass:        StateClassMeasurement,
			UnitOfMeasurement: UnitWatt,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Charging Power",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "charging_power"),
				Icon:        IconBatteryPlus,
				DeviceClass: DeviceClassPower,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.charge_w }}",
			},
			StateClass:        StateClassMeasurement,
			UnitOfMeasurement: UnitWatt,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Discharge Power",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "discharge_power"),
				Icon:        IconBatteryMinus,
				DeviceClass: DeviceClassPower,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.discharge_w }}",
			},
			StateClass:        StateClassMeasurement,
			UnitOfMeasurement: UnitWatt,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Generation Total",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "generation_total"),
				DeviceClass: DeviceClassEnergy,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.generation_total_kwh }}",
			},
			StateClass:        StateClassTotalIncreasing,
			UnitOfMeasurement: UnitKilowattHours,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Generation Today",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "generation_today"),
				DeviceClass: DeviceClassEnergy,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.generation_today_kwh }}",
			},
			StateClass:        StateClassTotalIncreasing,
			UnitOfMeasurement: UnitKilowattHours,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "SoC",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "soc"),
				DeviceClass: DeviceClassBattery,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.soc }}",
			},
			StateClass:        StateClassMeasurement,
			UnitOfMeasurement: UnitPercent,
		},
		{
			CommonConfig: CommonConfig{
				Name:     "Number Of Batteries",
				UniqueId: fmt.Sprintf("%s_%s", info.SerialNumber, "battery_num"),
				Icon:     IconCarBattery,
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.battery_num }}",
			},
			StateClass: StateClassMeasurement,
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Working Mode",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "work_mode"),
				DeviceClass: DeviceClassEnum,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.work_mode }}",
			},
			Options: []string{models.WorkModeLoadFirst, models.WorkModeBatteryFirst},
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Status",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "status"),
				DeviceClass: DeviceClassEnum,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.StateTopic,
				ValueTemplate: "{{ value_json.status }}",
			},
			Options: []string{
				models.Offline,
				models.WorkModeLoadFirst,
				models.WorkModeBatteryFirst,
				models.SmartSelfUse,
				models.Fault,
				models.Heating,
				models.OnGrid,
				models.OffGrid},
		},
	}

	for _, b := range info.Batteries {
		sensors = append(sensors, []Sensor{
			{
				CommonConfig: CommonConfig{
					Name:        fmt.Sprintf("%s SoC", b.Alias),
					UniqueId:    fmt.Sprintf("%s_%s_%s", info.SerialNumber, b.Alias, "soc"),
					DeviceClass: DeviceClassBattery,
					Device:      device,
					Origin:      origin,
				},
				StateConfig: StateConfig{
					StateTopic:    b.StateTopic,
					ValueTemplate: "{{ value_json.soc }}",
				},
				StateClass:        StateClassMeasurement,
				UnitOfMeasurement: UnitPercent,
			},
			{
				CommonConfig: CommonConfig{
					Name:        fmt.Sprintf("%s Temperature", b.Alias),
					UniqueId:    fmt.Sprintf("%s_%s_%s", info.SerialNumber, b.Alias, "temp"),
					DeviceClass: DeviceClassTemperature,
					Device:      device,
					Origin:      origin,
				},
				StateConfig: StateConfig{
					StateTopic:    b.StateTopic,
					ValueTemplate: "{{ value_json.temp }}",
				},
				StateClass:        StateClassMeasurement,
				UnitOfMeasurement: UnitCelsius,
			},
		}...)
	}

	return sensors
}
