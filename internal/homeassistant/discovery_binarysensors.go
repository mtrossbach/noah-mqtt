package homeassistant

import "fmt"

func generateBinarySensorDiscoveryPayload(appVersion string, info DeviceInfo) []BinarySensor {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	binarySensors := []BinarySensor{
		{
			Name:          "Connectivity",
			Icon:          "",
			DeviceClass:   DeviceClassConnectivity,
			ValueTemplate: "{{ value_json.status }}",
			PayloadOff:    "offline",
			PayloadOn:     "online",
			UniqueId:      fmt.Sprintf("%s_connectivity", info.SerialNumber),
			StateTopic:    info.StateTopic,
			Device:        device,
			Origin:        origin,
		},
	}

	return binarySensors
}
