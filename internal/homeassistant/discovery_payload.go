package homeassistant

import "fmt"

func generateSensorDiscoveryPayload(swVersion string, info DeviceInfo) []Sensor {
	device := Device{
		Identifiers:  []string{fmt.Sprintf("noah_%s", info.SerialNumber)},
		Name:         info.Alias,
		Manufacturer: "Growatt",
		HwVersion:    info.Version,
		SwVersion:    swVersion,
		Model:        info.Model,
		SerialNumber: info.SerialNumber,
	}

	sensors := []Sensor{
		{
			Name:              "Output Power",
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.output_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "output_power"),
			Device:            device,
		},
		{
			Name:              "Solar Power",
			Icon:              IconSolarPower,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.solar_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "solar_power"),
			Device:            device,
		},
		{
			Name:              "Charging Power",
			Icon:              IconBatteryPlus,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.charge_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "charging_power"),
			Device:            device,
		},
		{
			Name:              "Discharge Power",
			Icon:              IconBatteryMinus,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.discharge_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "discharge_power"),
			Device:            device,
		},
		{
			Name:              "Generation Total",
			DeviceClass:       DeviceClassEnergy,
			StateClass:        StateClassTotalIncreasing,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitKilowattHours,
			ValueTemplate:     "{{ value_json.generation_total_kwh }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "generation_total"),
			Device:            device,
		},
		{
			Name:              "Generation Today",
			DeviceClass:       DeviceClassEnergy,
			StateClass:        StateClassTotalIncreasing,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitKilowattHours,
			ValueTemplate:     "{{ value_json.generation_today_kwh }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "generation_today"),
			Device:            device,
		},
		{
			Name:              "SoC",
			DeviceClass:       DeviceClassBattery,
			StateClass:        StateClassMeasurement,
			StateTopic:        info.StateTopic,
			UnitOfMeasurement: UnitPercent,
			ValueTemplate:     "{{ value_json.soc }}",
			UniqueId:          fmt.Sprintf("%s_%s", info.SerialNumber, "soc"),
			Device:            device,
		},
		{
			Name:          "Number Of Batteries",
			StateClass:    StateClassMeasurement,
			StateTopic:    info.StateTopic,
			Icon:          IconCarBattery,
			ValueTemplate: "{{ value_json.battery_num }}",
			UniqueId:      fmt.Sprintf("%s_%s", info.SerialNumber, "battery_num"),
			Device:        device,
		},
	}

	for _, b := range info.Batteries {
		sensors = append(sensors, []Sensor{
			{
				Name:              fmt.Sprintf("%s SoC", b.Alias),
				DeviceClass:       DeviceClassBattery,
				StateClass:        StateClassMeasurement,
				StateTopic:        b.StateTopic,
				UnitOfMeasurement: UnitPercent,
				ValueTemplate:     "{{ value_json.soc }}",
				UniqueId:          fmt.Sprintf("%s_%s_%s", info.SerialNumber, b.Alias, "soc"),
				Device:            device,
			},
			{
				Name:              fmt.Sprintf("%s Temperature", b.Alias),
				DeviceClass:       DeviceClassTemperature,
				StateClass:        StateClassMeasurement,
				StateTopic:        b.StateTopic,
				UnitOfMeasurement: UnitCelsius,
				ValueTemplate:     "{{ value_json.temp }}",
				UniqueId:          fmt.Sprintf("%s_%s_%s", info.SerialNumber, b.Alias, "temp"),
				Device:            device,
			},
		}...)
	}

	return sensors
}
