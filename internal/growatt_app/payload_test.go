package growatt_app

import (
	"nexa-mqtt/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_devicePayload(t *testing.T) {
	noahStatus := NoahStatus{
		ResponseContainerV2[NoahStatusObj]{
			Obj: NoahStatusObj{
				LoadPower:     "400",
				GridPower:     "0",
				ChargePower:   "132",
				GroplugPower:  "0",
				WorkMode:      "0",
				Soc:           "93",
				EastronStatus: "-1",
				//AssociatedInvSn: nil,
				BatteryNum:     "1",
				ProfitToday:    "0",
				PlantID:        "10421077",
				DisChargePower: "0",
				EacTotal:       "9.6",
				EacToday:       "3.3",
				IsHaveCt:       "false",
				OnOffGrid:      "0",
				Pac:            "-400",
				Ppv:            "538",
				Alias:          "NEXA 2000",
				ProfitTotal:    "0",
				MoneyUnit:      "â¬",
				GroplugNum:     "0",
				OtherPower:     "-400",
				Status:         "6",
			},
		},
	}

	dp := devicePayload(&noahStatus)

	assert.Equal(t, -400.0, dp.OutputPower)
	assert.Equal(t, 538.0, dp.SolarPower)
	assert.Equal(t, 93.0, dp.Soc)
	assert.Equal(t, 132.0, dp.ChargePower)
	assert.Equal(t, 0.0, dp.DischargePower)
	assert.Equal(t, 1, dp.BatteryNum)
	assert.Equal(t, 9.6, dp.GenerationTotalEnergy)
	assert.Equal(t, 3.3, dp.GenerationTodayEnergy)
	assert.Equal(t, models.WorkMode("load_first"), dp.WorkMode)
	assert.Equal(t, "on_grid", dp.Status)
}

func Test_batteryPayload(t *testing.T) {
	batteryDetails := BatteryDetails{
		Temp:      "39",
		SerialNum: "serial123",
		Soc:       "93",
	}

	bp := batteryPayload(&batteryDetails)

	assert.Equal(t, "serial123", bp.SerialNumber)
	assert.Equal(t, 93.0, bp.Soc)
	assert.Equal(t, 39.0, bp.Temperature)
}

func Test_parameterPayload(t *testing.T) {
	nexaInfo := NexaInfo{}

	nexaInfo.Obj.Noah.ChargingSocHighLimit = "95"
	nexaInfo.Obj.Noah.DefaultMode = "0"
	nexaInfo.Obj.Noah.DefaultACCouplePower = "100"
	nexaInfo.Obj.Noah.ChargingSocLowLimit = "11"

	pp := parameterPayload(&nexaInfo)

	assert.Equal(t, 95.0, pp.ChargingLimit)
}
