package homeassistant

import "fmt"

func generateDevice(info DeviceInfo) Device {
	return Device{
		Identifiers:  []string{fmt.Sprintf("nexa_%s", info.SerialNumber)},
		Name:         info.Alias,
		Manufacturer: "Growatt",
		SwVersion:    info.Version,
		Model:        info.Model,
		SerialNumber: info.SerialNumber,
	}
}

func generateOrigin(appVersion string) Origin {
	return Origin{
		Name:       "nexa-mqtt",
		SwVersion:  appVersion,
		SupportUrl: "https://github.com/mgerczuk/nexa-mqtt",
	}
}
