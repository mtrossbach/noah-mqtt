package growatt

import (
	"noah-mqtt/pkg/models"
	"strconv"
)

type LoginResult struct {
	Back struct {
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
		User    struct {
			ID int `json:"id"`
		} `json:"user"`
	} `json:"back"`
}

type PlantList struct {
	Back struct {
		Data []struct {
			PlantID string `json:"plantId"`
		} `json:"data"`
		Success bool `json:"success"`
	} `json:"back"`
}

type NoahPlantInfo struct {
	Msg    string `json:"msg"`
	Result int    `json:"result"`
	Obj    struct {
		IsPlantNoahSystem bool   `json:"isPlantNoahSystem"`
		PlantID           string `json:"plantId"`
		IsPlantHaveNoah   bool   `json:"isPlantHaveNoah"`
		DeviceSn          string `json:"deviceSn"`
		PlantName         string `json:"plantName"`
	} `json:"obj"`
}

type NoahStatus struct {
	Msg    string `json:"msg"`
	Result int    `json:"result"`
	Obj    struct {
		ChargePower     string `json:"chargePower"`
		WorkMode        string `json:"workMode"`
		Soc             string `json:"soc"`
		AssociatedInvSn string `json:"associatedInvSn"`
		BatteryNum      string `json:"batteryNum"`
		ProfitToday     string `json:"profitToday"`
		PlantID         string `json:"plantId"`
		DisChargePower  string `json:"disChargePower"`
		EacTotal        string `json:"eacTotal"`
		EacToday        string `json:"eacToday"`
		Pac             string `json:"pac"`
		Ppv             string `json:"ppv"`
		Alias           string `json:"alias"`
		ProfitTotal     string `json:"profitTotal"`
		MoneyUnit       string `json:"moneyUnit"`
		Status          string `json:"status"`
	} `json:"obj"`
}

func (n *NoahStatus) ToPayload() models.Payload {
	return models.Payload{
		OutputPower:           parseFloat(n.Obj.Pac),
		SolarPower:            parseFloat(n.Obj.Ppv),
		Soc:                   parseFloat(n.Obj.Soc),
		ChargePower:           parseFloat(n.Obj.ChargePower),
		DischargePower:        parseFloat(n.Obj.DisChargePower),
		BatteryCount:          int(parseFloat(n.Obj.BatteryNum)),
		GenerationTotalEnergy: parseFloat(n.Obj.EacTotal),
		GenerationTodayEnergy: parseFloat(n.Obj.EacToday),
		WorkMode:              workModeFromString(n.Obj.WorkMode),
	}
}

func workModeFromString(s string) models.WorkMode {
	if s == "0" {
		return models.WorkModeLoadFirst
	}
	return models.WorkModeBatteryFirst
}

func parseFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	} else {
		return 0
	}
}
