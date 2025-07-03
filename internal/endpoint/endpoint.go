package endpoint

import "nexa-mqtt/pkg/models"

type Endpoint interface {
	SetParameterApplier(applier ParameterApplier)
	SetDevices(devices []models.NoahDevicePayload)
	PublishDeviceStatus(device models.NoahDevicePayload, status models.DevicePayload)
	PublishBatteryDetails(device models.NoahDevicePayload, details []models.BatteryPayload)
	PublishParameterData(device models.NoahDevicePayload, param models.ParameterPayload)
}
