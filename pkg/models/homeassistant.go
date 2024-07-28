package models

type HomeAssistantSensor struct {
	Name              string              `json:"name"`
	Icon              string              `json:"icon,omitempty"`
	DeviceClass       string              `json:"device_class"`
	StateTopic        string              `json:"state_topic"`
	UnitOfMeasurement string              `json:"unit_of_measurement"`
	ValueTemplate     string              `json:"value_template"`
	UniqueId          string              `json:"unique_id"`
	Device            HomeAssistantDevice `json:"device"`
}

type HomeAssistantDevice struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
	SerialNumber string   `json:"serial_number"`
}
