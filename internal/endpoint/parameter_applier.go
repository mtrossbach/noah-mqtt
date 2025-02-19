package endpoint

import "noah-mqtt/pkg/models"

type ParameterApplier interface {
	SetOutputPowerW(device models.NoahDevicePayload, power float64) bool
	SetChargingLimit(device models.NoahDevicePayload, limit float64) bool
	SetDischargeLimit(device models.NoahDevicePayload, limit float64) bool
}
