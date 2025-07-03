package homeassistant

type DeviceClass string

const (
	DeviceClassNone         DeviceClass = ""
	DeviceClassEnergy       DeviceClass = "energy"
	DeviceClassBattery      DeviceClass = "battery"
	DeviceClassTemperature  DeviceClass = "temperature"
	DeviceClassPower        DeviceClass = "power"
	DeviceClassConnectivity DeviceClass = "connectivity"
	DeviceClassEnum         DeviceClass = "enum"
)

type StateClass string

const (
	StateClassMeasurement     StateClass = "measurement"
	StateClassTotalIncreasing StateClass = "total_increasing"
)

type Unit string

const (
	UnitKilowattHours Unit = "kWh"
	UnitWatt          Unit = "W"
	UnitPercent       Unit = "%"
	UnitCelsius       Unit = "°C"
)

type Icon string

const (
	IconSolarPower              Icon = "mdi:solar-power"
	IconBatteryPlus             Icon = "mdi:battery-plus"
	IconBatteryMinus            Icon = "mdi:battery-minus"
	IconCarBattery              Icon = "mdi:car-battery"
	IconHeatWave                Icon = "mdi:heat-wave"
	IconBatteryArrowUpOutline   Icon = "mdi:battery-arrow-up-outline"
	IconBatteryArrowDownOutline Icon = "mdi:battery-arrow-down-outline"
)

type Device struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Name         string   `json:"name,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	HwVersion    string   `json:"hw_version,omitempty"`
	SwVersion    string   `json:"sw_version,omitempty"`
	Model        string   `json:"model,omitempty"`
	SerialNumber string   `json:"serial_number,omitempty"`
}

type Origin struct {
	Name       string `json:"name"`
	SwVersion  string `json:"sw_version,omitempty"`
	SupportUrl string `json:"support_url,omitempty"`
}

type Mode string

const (
	ModeBox    Mode = "box"
	ModeSlider Mode = "slider"
)

type CommonConfig struct {
	Name        string      `json:"name"`
	UniqueId    string      `json:"unique_id,omitempty"`
	Icon        Icon        `json:"icon,omitempty"`
	DeviceClass DeviceClass `json:"device_class,omitempty"`
	Device      Device      `json:"device,omitempty"`
	Origin      Origin      `json:"origin,omitempty"`
}

type StateConfig struct {
	StateTopic    string `json:"state_topic"`
	ValueTemplate string `json:"value_template,omitempty"`
}

type CommandConfig struct {
	CommandTopic    string `json:"command_topic,omitempty"`
	CommandTemplate string `json:"command_template,omitempty"`
}

type BinarySensor struct {
	CommonConfig
	StateConfig
	PayloadOff string `json:"payload_off,omitempty"`
	PayloadOn  string `json:"payload_on,omitempty"`
}

type Sensor struct {
	CommonConfig
	StateConfig
	StateClass        StateClass `json:"state_class,omitempty"`
	UnitOfMeasurement Unit       `json:"unit_of_measurement,omitempty"`
	Options           []string   `json:"options,omitempty"`
}

type Select struct {
	CommonConfig
	StateConfig
	CommandConfig
	StateClass        StateClass `json:"state_class,omitempty"`
	UnitOfMeasurement Unit       `json:"unit_of_measurement,omitempty"`
	Options           []string   `json:"options,omitempty"`
	Component         string     `json:"component,omitempty"`
}

type Number struct {
	CommonConfig
	StateConfig
	CommandConfig
	StateClass        StateClass `json:"state_class,omitempty"`
	UnitOfMeasurement Unit       `json:"unit_of_measurement,omitempty"`
	Mode              Mode       `json:"mode,omitempty"`
	Step              float64    `json:"step,omitempty"`
	Min               float64    `json:"min,omitempty"`
	Max               float64    `json:"max,omitempty"`
}
