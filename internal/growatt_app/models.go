package growatt_app

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

type NoahPlantInfoObj struct {
	IsPlantNoahSystem bool   `json:"isPlantNoahSystem"`
	PlantID           string `json:"plantId"`
	IsPlantHaveNoah   bool   `json:"isPlantHaveNoah"`
	IsPlantHaveNexa   bool   `json:"isPlantHaveNexa"`
	DeviceSn          string `json:"deviceSn"`
	PlantName         string `json:"plantName"`
}

type NoahPlantInfo struct {
	ResponseContainerV2[NoahPlantInfoObj]
}

type NoahStatusObj struct {
	LoadPower       string `json:"loadPower"` // new
	GridPower       string `json:"gridPower"` // new
	ChargePower     string `json:"chargePower"`
	GroplugPower    string `json:"groplugPower"` // new
	WorkMode        string `json:"workMode"`
	Soc             string `json:"soc"`
	EastronStatus   string `json:"eastronStatus"` // new
	AssociatedInvSn string `json:"associatedInvSn"`
	BatteryNum      string `json:"batteryNum"`
	ProfitToday     string `json:"profitToday"`
	PlantID         string `json:"plantId"`
	DisChargePower  string `json:"disChargePower"`
	EacTotal        string `json:"eacTotal"`
	EacToday        string `json:"eacToday"`
	IsHaveCt        string `json:"isHaveCt"`  // new
	OnOffGrid       string `json:"onOffGrid"` // new
	Pac             string `json:"pac"`
	Ppv             string `json:"ppv"`
	Alias           string `json:"alias"`
	ProfitTotal     string `json:"profitTotal"`
	MoneyUnit       string `json:"moneyUnit"`
	GroplugNum      string `json:"groplugNum"` // new
	OtherPower      string `json:"otherPower"` // new
	Status          string `json:"status"`     // 1 = online, -1 = offline, 5 = heating
}

type NoahStatus struct {
	ResponseContainerV2[NoahStatusObj]
}

type NexaInfoObj struct {
	Noah struct {
		TimeSegment                 []map[string]string `json:"time_segment"`
		AntiBackflowEnable          string              `json:"antiBackflowEnable"`          // new
		AcCouplePowerControl        string              `json:"acCouplePowerControl"`        // new
		AmmeterModel                string              `json:"ammeterModel"`                // new
		AmmeterSn                   string              `json:"ammeterSn"`                   // new
		ShellyList                  []interface{}       `json:"shellyList"`                  // new
		GridSet                     string              `json:"gridSet"`                     // new
		AntiBackflowPowerPercentage string              `json:"antiBackflowPowerPercentage"` // new
		BatSns                      []string            `json:"batSns"`
		ManName                     string              `json:"manName"`
		AssociatedInvSn             string              `json:"associatedInvSn"`
		PlantID                     string              `json:"plantId"`
		ChargingSocHighLimit        string              `json:"chargingSocHighLimit"`
		DefaultMode                 string              `json:"defaultMode"`          // new
		DefaultACCouplePower        string              `json:"defaultACCouplePower"` // new
		Version                     string              `json:"version"`
		DeviceSn                    string              `json:"deviceSn"`
		ChargingSocLowLimit         string              `json:"chargingSocLowLimit"`
		FormulaMoney                string              `json:"formulaMoney"`
		Alias                       string              `json:"alias"`
		Model                       string              `json:"model"`
		CtType                      string              `json:"ctType"`                // new
		AllowGridCharging           string              `json:"allowGridCharging"`     // new
		GridConnectionControl       string              `json:"gridConnectionControl"` // new
		PlantName                   string              `json:"plantName"`
		AssociatedInvManAndModel    int                 `json:"associatedInvManAndModel"`
		TempType                    string              `json:"tempType"`
		MoneyUnitText               string              `json:"moneyUnitText"`
	} `json:"noah"`
	PlantList []struct {
		PlantID      string      `json:"plantId"`
		PlantImgName interface{} `json:"plantImgName"`
		PlantName    string      `json:"plantName"`
	} `json:"plantList"`
	UnitList map[string]string `json:"unitList"` // new
}

type NexaInfo struct {
	ResponseContainerV2[NexaInfoObj]
}

type BatteryInfoObj struct {
	Batter   []BatteryDetails `json:"batter"`
	TempType string           `json:"tempType"`
	Time     string           `json:"time"`
}

type BatteryInfo struct {
	ResponseContainerV2[BatteryInfoObj]
}

type BatteryDetails struct {
	SerialNum string `json:"serialNum"`
	Soc       string `json:"soc"`
	Temp      string `json:"temp"`
}

type SetResponse struct {
	ResponseContainerV2[any]
}
