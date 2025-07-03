package homeassistant

import "fmt"

func generateNumberDiscoveryPayload(appVersion string, info DeviceInfo) []Number {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	numbers := []Number{
		{
			CommonConfig: CommonConfig{
				Name:        "Default AC Output Power",
				UniqueId:    fmt.Sprintf("%s_default_output_w", info.SerialNumber),
				Icon:        "",
				DeviceClass: DeviceClassPower,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.default_output_w }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"default_output_w\": {{ value }}}",
			},
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               0,
			Max:               800,
			UnitOfMeasurement: UnitWatt,
		},
		{
			CommonConfig: CommonConfig{
				Name:     "Charging Limit",
				UniqueId: fmt.Sprintf("%s_charging_limit", info.SerialNumber),
				Icon:     IconBatteryArrowUpOutline,
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.charging_limit }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"charging_limit\": {{ value }}}",
			},
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               70,
			Max:               100,
			UnitOfMeasurement: UnitPercent,
		},
		{
			CommonConfig: CommonConfig{
				Name:     "Discharge Limit",
				UniqueId: fmt.Sprintf("%s_discharge_limit", info.SerialNumber),
				Icon:     IconBatteryArrowDownOutline,
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.discharge_limit }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"discharge_limit\": {{ value }}}",
			},
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               0,
			Max:               30,
			UnitOfMeasurement: UnitPercent,
		},
	}

	return numbers
}
