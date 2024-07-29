package homeassistant

type DeviceClass string

const (
	DeviceClassEnergy  DeviceClass = "energy"
	DeviceClassBattery DeviceClass = "battery"
	DeviceClassPower   DeviceClass = "power"
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
)

type Icon string

const (
	IconSolarPower   Icon = "mdi:solar-power"
	IconBatteryPlus  Icon = "mdi:battery-plus"
	IconBatteryMinus Icon = "mdi:battery-minus"
)

type Sensor struct {
	Name              string      `json:"name"`
	Icon              Icon        `json:"icon,omitempty"`
	DeviceClass       DeviceClass `json:"device_class,omitempty"`
	StateTopic        string      `json:"state_topic"`
	StateClass        StateClass  `json:"state_class,omitempty"`
	UnitOfMeasurement Unit        `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string      `json:"value_template,omitempty"`
	UniqueId          string      `json:"unique_id,omitempty"`
	Device            Device      `json:"device,omitempty"`
}

type Device struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Name         string   `json:"name,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	SerialNumber string   `json:"serial_number,omitempty"`
}
