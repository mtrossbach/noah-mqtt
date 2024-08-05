package homeassistant

import "fmt"

func generateNumberDiscoveryPayload(appVersion string, info DeviceInfo) []Number {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	numbers := []Number{
		{
			Name:              "System Output Power",
			UniqueId:          fmt.Sprintf("%s_system_output_power", info.SerialNumber),
			CommandTemplate:   "{\"output_power_w\": {{ value }}}",
			CommandTopic:      info.ParameterCommandTopic,
			Device:            device,
			Origin:            origin,
			Icon:              "",
			DeviceClass:       DeviceClassPower,
			StateTopic:        info.ParameterStateTopic,
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               0,
			Max:               800,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.output_power_w }}",
		},
	}

	return numbers
}
