package growatt

type TokenResponse struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Token string `json:"token"`
}

type LoginResult struct {
	Back struct {
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
		User    struct {
			ID int `json:"id"`
		} `json:"user"`
	} `json:"back"`
}

type PlantListV2 struct {
	PlantList []struct {
		ID int `json:"id"`
	} `json:"PlantList"`
}

type ResponseContainerV2[T any] struct {
	Msg    string `json:"msg"`
	Result int    `json:"result"`
	Obj    T      `json:"obj"`
}

type NoahPlantInfo struct {
	ResponseContainerV2[struct {
		IsPlantNoahSystem bool   `json:"isPlantNoahSystem"`
		PlantID           string `json:"plantId"`
		IsPlantHaveNoah   bool   `json:"isPlantHaveNoah"`
		DeviceSn          string `json:"deviceSn"`
		PlantName         string `json:"plantName"`
	}]
}

type NoahStatus struct {
	ResponseContainerV2[struct {
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
		Status          string `json:"status"` // 1 = online, -1 = offline
	}]
}

type NoahInfo struct {
	ResponseContainerV2[struct {
		Noah struct {
			TimeSegment              []map[string]string `json:"time_segment"`
			BatSns                   []string            `json:"batSns"`
			ManName                  string              `json:"manName"`
			AssociatedInvSn          string              `json:"associatedInvSn"`
			PlantID                  string              `json:"plantId"`
			ChargingSocHighLimit     string              `json:"chargingSocHighLimit"`
			DefaultPower             string              `json:"defaultPower"`
			Version                  string              `json:"version"`
			DeviceSn                 string              `json:"deviceSn"`
			ChargingSocLowLimit      string              `json:"chargingSocLowLimit"`
			FormulaMoney             string              `json:"formulaMoney"`
			ModelName                string              `json:"modelName"`
			Alias                    string              `json:"alias"`
			Model                    string              `json:"model"`
			PlantName                string              `json:"plantName"`
			AssociatedInvManAndModel int                 `json:"associatedInvManAndModel"`
			TempType                 string              `json:"tempType"`
			MoneyUnitText            string              `json:"moneyUnitText"`
		} `json:"noah"`
		PlantList []struct {
			PlantID      string      `json:"plantId"`
			PlantImgName interface{} `json:"plantImgName"`
			PlantName    string      `json:"plantName"`
		} `json:"plantList"`
	}]
}

type BatteryInfo struct {
	ResponseContainerV2[struct {
		Batter   []BatteryDetails `json:"batter"`
		TempType string           `json:"tempType"`
		Time     string           `json:"time"`
	}]
}

type BatteryDetails struct {
	SerialNum string `json:"serialNum"`
	Soc       string `json:"soc"`
	Temp      string `json:"temp"`
}
