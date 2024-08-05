package homeassistant

import "fmt"

func generateDevice(info DeviceInfo) Device {
	return Device{
		Identifiers:  []string{fmt.Sprintf("noah_%s", info.SerialNumber)},
		Name:         info.Alias,
		Manufacturer: "Growatt",
		SwVersion:    info.Version,
		Model:        info.Model,
		SerialNumber: info.SerialNumber,
	}
}

func generateOrigin(appVersion string) Origin {
	return Origin{
		Name:       "noah-mqtt",
		SwVersion:  appVersion,
		SupportUrl: "https://github.com/mtrossbach/noah-mqtt",
	}
}
