package homeassistant

type DeviceInfo struct {
	SerialNumber string
	Model        string
	Version      string
	Alias        string
	StateTopic   string
	Batteries    []BatteryInfo
}

type BatteryInfo struct {
	Alias      string
	StateTopic string
}
