package homeassistant

import "fmt"

func generateSensorDiscoveryPayload(deviceName string, serialNumber string, batteries []BatteryInfo, stateTopic string) []Sensor {
	device := Device{
		Identifiers:  []string{fmt.Sprintf("noah_%s", serialNumber)},
		Name:         deviceName,
		Manufacturer: "Growatt",
		Model:        "Noah",
		SerialNumber: serialNumber,
	}

	sensors := []Sensor{
		{
			Name:              "Output Power",
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.output_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "output_power"),
			Device:            device,
		},
		{
			Name:              "Solar Power",
			Icon:              IconSolarPower,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.solar_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "solar_power"),
			Device:            device,
		},
		{
			Name:              "Charging Power",
			Icon:              IconBatteryPlus,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.charge_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "charging_power"),
			Device:            device,
		},
		{
			Name:              "Discharge Power",
			Icon:              IconBatteryMinus,
			DeviceClass:       DeviceClassPower,
			StateClass:        StateClassMeasurement,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.discharge_w }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "discharge_power"),
			Device:            device,
		},
		{
			Name:              "Generation Total",
			DeviceClass:       DeviceClassEnergy,
			StateClass:        StateClassTotalIncreasing,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitKilowattHours,
			ValueTemplate:     "{{ value_json.generation_total_kwh }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "generation_total"),
			Device:            device,
		},
		{
			Name:              "Generation Today",
			DeviceClass:       DeviceClassEnergy,
			StateClass:        StateClassTotalIncreasing,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitKilowattHours,
			ValueTemplate:     "{{ value_json.generation_today_kwh }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "generation_today"),
			Device:            device,
		},
		{
			Name:              "SoC",
			DeviceClass:       DeviceClassBattery,
			StateClass:        StateClassMeasurement,
			StateTopic:        stateTopic,
			UnitOfMeasurement: UnitPercent,
			ValueTemplate:     "{{ value_json.soc }}",
			UniqueId:          fmt.Sprintf("%s_%s", serialNumber, "soc"),
			Device:            device,
		},
		{
			Name:          "Number Of Batteries",
			StateClass:    StateClassMeasurement,
			StateTopic:    stateTopic,
			Icon:          IconCarBattery,
			ValueTemplate: "{{ value_json.battery_num }}",
			UniqueId:      fmt.Sprintf("%s_%s", serialNumber, "battery_num"),
			Device:        device,
		},
	}

	for _, b := range batteries {
		sensors = append(sensors, []Sensor{
			{
				Name:              fmt.Sprintf("%s SoC", b.Alias),
				DeviceClass:       DeviceClassBattery,
				StateClass:        StateClassMeasurement,
				StateTopic:        b.StateTopic,
				UnitOfMeasurement: UnitPercent,
				ValueTemplate:     "{{ value_json.soc }}",
				UniqueId:          fmt.Sprintf("%s_%s_%s", serialNumber, b.Alias, "soc"),
				Device:            device,
			},
			{
				Name:              fmt.Sprintf("%s Temperature", b.Alias),
				DeviceClass:       DeviceClassTemperature,
				StateClass:        StateClassMeasurement,
				StateTopic:        b.StateTopic,
				UnitOfMeasurement: UnitCelsius,
				ValueTemplate:     "{{ value_json.temp }}",
				UniqueId:          fmt.Sprintf("%s_%s_%s", serialNumber, b.Alias, "temp"),
				Device:            device,
			},
		}...)
	}

	return sensors
}
