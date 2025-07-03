package endpoint

import "nexa-mqtt/pkg/models"

type ParameterApplier interface {
	SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) bool
	SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) bool
}
