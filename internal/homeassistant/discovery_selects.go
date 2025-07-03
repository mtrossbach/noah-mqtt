package homeassistant

import (
	"fmt"
	"nexa-mqtt/pkg/models"
)

func generateSelectDiscoveryPayload(appVersion string, info DeviceInfo) []Select {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	selects := []Select{
		{
			CommonConfig: CommonConfig{
				Name:        "Default Mode",
				UniqueId:    fmt.Sprintf("%s_%s", info.SerialNumber, "default_mode"),
				DeviceClass: DeviceClassEnum,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.default_mode }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"default_mode\": \"{{ value }}\"}",
			},
			Options:   []string{models.WorkModeLoadFirst, models.WorkModeBatteryFirst},
			Component: "select",
		},
	}

	return selects
}
