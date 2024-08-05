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
		FormulaCo2         int           `json:"formulaCo2"`
		CompanyName        string        `json:"companyName"`
		EtodayCo2Text      string        `json:"etodayCo2Text"`
		UserBean           interface{}   `json:"userBean"`
		FormulaSo2         int           `json:"formulaSo2"`
		GridPort           string        `json:"gridPort"`
		Children           []interface{} `json:"children"`
		PlantFromBean      interface{}   `json:"plantFromBean"`
		ID                 int           `json:"id"`
		EYearMoneyText     string        `json:"EYearMoneyText"`
		TempType           int           `json:"tempType"`
		EtotalCoalText     string        `json:"etotalCoalText"`
		EtotalSo2Text      string        `json:"etotalSo2Text"`
		PlantLng           string        `json:"plant_lng"`
		LocationImgName    string        `json:"locationImgName"`
		DeviceCount        int           `json:"deviceCount"`
		MapCountryID       int           `json:"map_countryId"`
		MapLat             string        `json:"mapLat"`
		PrMonth            string        `json:"prMonth"`
		EtotalMoney        int           `json:"etotalMoney"`
		PlantType          int           `json:"plantType"`
		WindAngle          int           `json:"windAngle"`
		FormulaMoney       int           `json:"formulaMoney"`
		MapCity            string        `json:"mapCity"`
		NominalPower       int           `json:"nominalPower"`
		LogoImgName        string        `json:"logoImgName"`
		LatitudeText       string        `json:"latitudeText"`
		UserAccount        string        `json:"userAccount"`
		StorageTodayToUser int           `json:"storage_TodayToUser"`
		EventMessBeanList  []interface{} `json:"eventMessBeanList"`
		MapCityID          int           `json:"map_cityId"`
		CreateDateTextA    string        `json:"createDateTextA"`
		Status             int           `json:"status"`
		FormulaMoneyUnitID string        `json:"formulaMoneyUnitId"`
		EnergyMonth        int           `json:"energyMonth"`
		City               string        `json:"city"`
		PrToday            string        `json:"prToday"`
		EtodayCoalText     string        `json:"etodayCoalText"`
		CurrentPac         int           `json:"currentPac"`
		ParentID           string        `json:"parentID"`
		PlantAddress       string        `json:"plantAddress"`
		EnvTemp            int           `json:"envTemp"`
		FormulaCoal        int           `json:"formulaCoal"`
		TreeID             string        `json:"treeID"`
		HasStorage         int           `json:"hasStorage"`
		StorageTotalToUser int           `json:"storage_TotalToUser"`
		FixedPowerPrice    int           `json:"fixedPowerPrice"`
		EtodaySo2Text      string        `json:"etodaySo2Text"`
		PanelTemp          int           `json:"panelTemp"`
		CreateDate         struct {
			Date           int   `json:"date"`
			Hours          int   `json:"hours"`
			Seconds        int   `json:"seconds"`
			Month          int   `json:"month"`
			TimezoneOffset int   `json:"timezoneOffset"`
			Year           int   `json:"year"`
			Minutes        int   `json:"minutes"`
			Time           int64 `json:"time"`
			Day            int   `json:"day"`
		} `json:"createDate"`
		MapProvinceID            int           `json:"map_provinceId"`
		PairViewUserAccount      string        `json:"pairViewUserAccount"`
		EmonthSo2Text            string        `json:"emonthSo2Text"`
		PeakPeriodPrice          int           `json:"peakPeriodPrice"`
		HasDeviceOnLine          int           `json:"hasDeviceOnLine"`
		StorageBattoryPercentage int           `json:"storage_BattoryPercentage"`
		EtodayMoney              int           `json:"etodayMoney"`
		FormulaTree              int           `json:"formulaTree"`
		PlantNmi                 string        `json:"plantNmi"`
		ProtocolID               string        `json:"protocolId"`
		GridServerURL            string        `json:"gridServerUrl"`
		MoneyUnitText            string        `json:"moneyUnitText"`
		LongitudeD               string        `json:"longitude_d"`
		Country                  string        `json:"country"`
		LongitudeF               string        `json:"longitude_f"`
		EtodayMoneyText          string        `json:"etodayMoneyText"`
		LongitudeM               string        `json:"longitude_m"`
		PhoneNum                 string        `json:"phoneNum"`
		StorageTodayToGrid       int           `json:"storage_TodayToGrid"`
		DesignCompany            string        `json:"designCompany"`
		InstallMapName           string        `json:"installMapName"`
		CurrentPacStr            string        `json:"currentPacStr"`
		EtotalMoneyText          string        `json:"etotalMoneyText"`
		WindSpeed                int           `json:"windSpeed"`
		ValleyPeriodPrice        int           `json:"valleyPeriodPrice"`
		LatitudeF                string        `json:"latitude_f"`
		MapLng                   string        `json:"mapLng"`
		LatitudeD                string        `json:"latitude_d"`
		Level                    int           `json:"level"`
		LatitudeM                string        `json:"latitude_m"`
		EnergyYear               int           `json:"energyYear"`
		LongitudeText            string        `json:"longitudeText"`
		FlatPeriodPrice          int           `json:"flatPeriodPrice"`
		EmonthCoalText           string        `json:"emonthCoalText"`
		ParamBean                interface{}   `json:"paramBean"`
		EtotalCo2Text            string        `json:"etotalCo2Text"`
		ImgPath                  string        `json:"imgPath"`
		IsShare                  bool          `json:"isShare"`
		PlantLat                 string        `json:"plant_lat"`
		EmonthCo2Text            string        `json:"emonthCo2Text"`
		Timezone                 int           `json:"timezone"`
		GridCompany              string        `json:"gridCompany"`
		StorageEChargeToday      int           `json:"storage_eChargeToday"`
		Remark                   string        `json:"remark"`
		StorageTotalToGrid       int           `json:"storage_TotalToGrid"`
		DefaultPlant             bool          `json:"defaultPlant"`
		CreateDateText           string        `json:"createDateText"`
		CurrentPacTxt            string        `json:"currentPacTxt"`
		UnitMap                  interface{}   `json:"unitMap"`
		AlarmValue               int           `json:"alarmValue"`
		TreeName                 string        `json:"treeName"`
		Alias                    string        `json:"alias"`
		Irradiance               int           `json:"irradiance"`
		FormulaMoneyStr          string        `json:"formulaMoneyStr"`
		OnLineEnvCount           int           `json:"onLineEnvCount"`
		StorageEDisChargeToday   int           `json:"storage_eDisChargeToday"`
		TimezoneText             string        `json:"timezoneText"`
		DataLogList              []interface{} `json:"dataLogList"`
		MapAreaID                int           `json:"map_areaId"`
		EtotalFormulaTreeText    string        `json:"etotalFormulaTreeText"`
		PlantImgName             string        `json:"plantImgName"`
		EToday                   float64       `json:"eToday"`
		ETotal                   float64       `json:"eTotal"`
		EmonthMoneyText          string        `json:"emonthMoneyText"`
		NominalPowerStr          string        `json:"nominalPowerStr"`
		PlantName                string        `json:"plantName"`
	} `json:"PlantList"`
	StatusMap struct {
		Offline   int `json:"offline"`
		FaultNum  int `json:"faultNum"`
		OnlineNum int `json:"onlineNum"`
		AllNum    int `json:"allNum"`
	} `json:"statusMap"`
	UsereTotalMoney  int     `json:"usereTotalMoney"`
	UsereTotal       float64 `json:"usereTotal"`
	PlantNum         int     `json:"plantNum"`
	CurrentPageNum   int     `json:"currentPageNum"`
	UsereTodayMoney  int     `json:"usereTodayMoney"`
	UsercurrentPac   int     `json:"usercurrentPac"`
	UsernominalPower int     `json:"usernominalPower"`
	TotalPageNum     int     `json:"totalPageNum"`
	MoneyUnitText    string  `json:"moneyUnitText"`
	UsereToday       float64 `json:"usereToday"`
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
		Status          string `json:"status"`
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
