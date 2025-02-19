package growatt_web

type GrowattResult struct {
	Result int `json:"result"`
}

type GrowattPlant struct {
	PlantId   string `json:"id"`
	PlantName string `json:"plantName"`
}

type GrowattPlantDevices struct {
	Result int `json:"result"`
	Obj    struct {
		CurrPage int `json:"currPage"`
		Pages    int `json:"pages"`
		PageSize int `json:"pageSize"`
		Count    int `json:"count"`
		Ind      int `json:"ind"`
		Datas    []struct {
			DeviceType      string `json:"deviceType"`
			PtoStatus       string `json:"ptoStatus"`
			ShellyDeviceSn  string `json:"shellyDeviceSn"`
			TimeServer      string `json:"timeServer"`
			AccountName     string `json:"accountName"`
			Timezone        string `json:"timezone"`
			PlantID         string `json:"plantId"`
			DeviceTypeName  string `json:"deviceTypeName"`
			NominalPower    string `json:"nominalPower"`
			BdcStatus       string `json:"bdcStatus"`
			EToday          string `json:"eToday"`
			EMonth          string `json:"eMonth"`
			DatalogTypeTest string `json:"datalogTypeTest"`
			ETotal          string `json:"eTotal"`
			Pac             string `json:"pac"`
			DatalogSn       string `json:"datalogSn"`
			Alias           string `json:"alias"`
			Location        string `json:"location"`
			DeviceModel     string `json:"deviceModel"`
			Sn              string `json:"sn"`
			PlantName       string `json:"plantName"`
			Status          string `json:"status"`
			LastUpdateTime  string `json:"lastUpdateTime"`
		} `json:"datas"`
		NotPager bool `json:"notPager"`
	} `json:"obj"`
}

type GrowattNoahList struct {
	CurrPage int `json:"currPage"`
	Pages    int `json:"pages"`
	PageSize int `json:"pageSize"`
	Count    int `json:"count"`
	Ind      int `json:"ind"`
	Datas    []struct {
		Time8Power                  string `json:"time8Power"`
		ShellyDeviceSn              string `json:"shellyDeviceSn"`
		Time5Mode                   string `json:"time5Mode"`
		Time7Enable                 string `json:"time7Enable"`
		Soc                         string `json:"soc"`
		Time4Start                  string `json:"time4Start"`
		Time2End                    string `json:"time2End"`
		ShellyFlag                  string `json:"shellyFlag"`
		Time6End                    string `json:"time6End"`
		Time9End                    string `json:"time9End"`
		Time6Mode                   string `json:"time6Mode"`
		Time4Enable                 string `json:"time4Enable"`
		Pac                         string `json:"pac"`
		Lost                        string `json:"lost"`
		Model                       string `json:"model"`
		Time7Power                  string `json:"time7Power"`
		DeviceType                  string `json:"deviceType"`
		Time3Start                  string `json:"time3Start"`
		Time6Start                  string `json:"time6Start"`
		Time1Power                  string `json:"time1Power"`
		Time7Mode                   string `json:"time7Mode"`
		Time6Enable                 string `json:"time6Enable"`
		Time1End                    string `json:"time1End"`
		ChargingSocHighLimit        string `json:"chargingSocHighLimit"`
		Time5End                    string `json:"time5End"`
		Time9Start                  string `json:"time9Start"`
		DefaultPower                string `json:"defaultPower"`
		Version                     string `json:"version"`
		Time3Power                  string `json:"time3Power"`
		ChargingSocLowLimit         string `json:"chargingSocLowLimit"`
		NominalPower                string `json:"nominalPower"`
		Time6Power                  string `json:"time6Power"`
		Time8Mode                   string `json:"time8Mode"`
		ComVersion                  string `json:"comVersion"`
		Time9Power                  string `json:"time9Power"`
		Time3Enable                 string `json:"time3Enable"`
		GridConnectionControl       string `json:"gridConnectionControl"`
		Status                      string `json:"status"`
		LastUpdateTime              string `json:"lastUpdateTime"`
		Time2Enable                 string `json:"time2Enable"`
		WorkMode                    string `json:"workMode"`
		AccountName                 string `json:"accountName"`
		Timezone                    string `json:"timezone"`
		AntiBackflowEnable          string `json:"antiBackflowEnable"`
		Time5Power                  string `json:"time5Power"`
		AcCouplePowerControl        string `json:"acCouplePowerControl"`
		Time7Start                  string `json:"time7Start"`
		Time9Mode                   string `json:"time9Mode"`
		Time4End                    string `json:"time4End"`
		Time1Start                  string `json:"time1Start"`
		Time7End                    string `json:"time7End"`
		EMonth                      string `json:"eMonth"`
		Dtc                         string `json:"dtc"`
		Time1Mode                   string `json:"time1Mode"`
		Time9Enable                 string `json:"time9Enable"`
		Alias                       string `json:"alias"`
		DatalogSn                   string `json:"datalogSn"`
		SysTime                     string `json:"sysTime"`
		FwVersion                   string `json:"fwVersion"`
		Sn                          string `json:"sn"`
		Time4Power                  string `json:"time4Power"`
		AntiBackflowPowerPercentage string `json:"antiBackflowPowerPercentage"`
		Time1Enable                 string `json:"time1Enable"`
		Address                     string `json:"address"`
		DatalogType                 string `json:"datalogType"`
		PlantID                     string `json:"plantId"`
		Time2Mode                   string `json:"time2Mode"`
		Time3End                    string `json:"time3End"`
		Time8End                    string `json:"time8End"`
		EToday                      string `json:"eToday"`
		Time8Start                  string `json:"time8Start"`
		ETotal                      string `json:"eTotal"`
		Time2Start                  string `json:"time2Start"`
		Time4Mode                   string `json:"time4Mode"`
		Time5Enable                 string `json:"time5Enable"`
		DefaultMode                 string `json:"defaultMode"`
		Ppv                         string `json:"ppv"`
		Time3Mode                   string `json:"time3Mode"`
		DeviceModel                 string `json:"deviceModel"`
		Time2Power                  string `json:"time2Power"`
		Time5Start                  string `json:"time5Start"`
		DefaultACCouplePower        string `json:"defaultACCouplePower"`
		PlantName                   string `json:"plantName"`
		Time8Enable                 string `json:"time8Enable"`
	} `json:"datas"`
	NotPager bool `json:"notPager"`
}

type GrowattNoahHistory struct {
	Result int `json:"result"`
	Obj    struct {
		EndDate string `json:"endDate"`
		Datas   []struct {
			FunctionCode                   int     `json:"functionCode"`
			DeviceSn                       string  `json:"deviceSn"`
			DatalogSn                      string  `json:"datalogSn"`
			Time                           string  `json:"time"`
			IsAgain                        int     `json:"isAgain"`
			Dtc                            int     `json:"dtc"`
			Status                         int     `json:"status"`
			MpptProtectStatus              int     `json:"mpptProtectStatus"`
			PdWarnStatus                   int     `json:"pdWarnStatus"`
			Pac                            float64 `json:"pac"`
			EacToday                       float64 `json:"eacToday"`
			EacMonth                       float64 `json:"eacMonth"`
			EacYear                        float64 `json:"eacYear"`
			EacTotal                       float64 `json:"eacTotal"`
			Ppv                            float64 `json:"ppv"`
			WorkMode                       int     `json:"workMode"`
			TotalBatteryPackChargingStatus int     `json:"totalBatteryPackChargingStatus"`
			TotalBatteryPackChargingPower  int     `json:"totalBatteryPackChargingPower"`
			BatteryPackageQuantity         int     `json:"batteryPackageQuantity"`
			TotalBatteryPackSoc            int     `json:"totalBatteryPackSoc"`
			HeatingStatus                  int     `json:"heatingStatus"`
			FaultStatus                    int     `json:"faultStatus"`
			Battery1SerialNum              string  `json:"battery1SerialNum"`
			Battery1Soc                    int     `json:"battery1Soc"`
			Battery1Temp                   float64 `json:"battery1Temp"`
			Battery1WarnStatus             int     `json:"battery1WarnStatus"`
			Battery1ProtectStatus          int     `json:"battery1ProtectStatus"`
			Battery2SerialNum              string  `json:"battery2SerialNum"`
			Battery2Soc                    int     `json:"battery2Soc"`
			Battery2Temp                   float64 `json:"battery2Temp"`
			Battery2WarnStatus             int     `json:"battery2WarnStatus"`
			Battery2ProtectStatus          int     `json:"battery2ProtectStatus"`
			Battery3SerialNum              string  `json:"battery3SerialNum"`
			Battery3Soc                    int     `json:"battery3Soc"`
			Battery3Temp                   float64 `json:"battery3Temp"`
			Battery3WarnStatus             int     `json:"battery3WarnStatus"`
			Battery3ProtectStatus          int     `json:"battery3ProtectStatus"`
			Battery4SerialNum              string  `json:"battery4SerialNum"`
			Battery4Soc                    int     `json:"battery4Soc"`
			Battery4Temp                   float64 `json:"battery4Temp"`
			Battery4WarnStatus             int     `json:"battery4WarnStatus"`
			Battery4ProtectStatus          int     `json:"battery4ProtectStatus"`
			SettableTimePeriod             int     `json:"settableTimePeriod"`
			AcCoupleWarnStatus             int     `json:"acCoupleWarnStatus"`
			AcCoupleProtectStatus          int     `json:"acCoupleProtectStatus"`
			CtFlag                         int     `json:"ctFlag"`
			TotalHouseholdLoad             float64 `json:"totalHouseholdLoad"`
			HouseholdLoadApartFromGroplug  float64 `json:"householdLoadApartFromGroplug"`
			OnOffGrid                      int     `json:"onOffGrid"`
			CtSelfPower                    float64 `json:"ctSelfPower"`
			PlugPower                      float64 `json:"plugPower"`
			SocketPower                    int     `json:"socketPower"`
			ChargeSocLimit                 int     `json:"chargeSocLimit"`
			DischargeSocLimit              int     `json:"dischargeSocLimit"`
			Pv1Voltage                     float64 `json:"pv1Voltage"`
			Pv1Current                     float64 `json:"pv1Current"`
			Pv1Temp                        float64 `json:"pv1Temp"`
			Pv2Voltage                     float64 `json:"pv2Voltage"`
			Pv2Current                     float64 `json:"pv2Current"`
			Pv2Temp                        float64 `json:"pv2Temp"`
			SystemTemp                     float64 `json:"systemTemp"`
			MaxCellVoltage                 float64 `json:"maxCellVoltage"`
			MinCellVoltage                 float64 `json:"minCellVoltage"`
			BatteryCycles                  int     `json:"batteryCycles"`
			BatterySoh                     int     `json:"batterySoh"`
			Pv3Voltage                     float64 `json:"pv3Voltage"`
			Pv3Current                     float64 `json:"pv3Current"`
			Pv3Temp                        float64 `json:"pv3Temp"`
			Pv4Voltage                     float64 `json:"pv4Voltage"`
			Pv4Current                     float64 `json:"pv4Current"`
			Pv4Temp                        float64 `json:"pv4Temp"`
			TotalReactivePower             float64 `json:"totalReactivePower"`
			TotalPowerFactor               float64 `json:"totalPowerFactor"`
			ForwardReactiveEnergy          float64 `json:"forwardReactiveEnergy"`
			ReverseReactiveEnergy          float64 `json:"reverseReactiveEnergy"`
			ApparentEnergy                 float64 `json:"apparentEnergy"`
			TotalActiveEnergy              float64 `json:"totalActiveEnergy"`
			TotalReactiveEnergy            float64 `json:"totalReactiveEnergy"`
			Frequency                      float64 `json:"frequency"`
			PhaseAReactivePower            float64 `json:"phaseAReactivePower"`
			PhaseBReactivePower            float64 `json:"phaseBReactivePower"`
			PhaseCReactivePower            float64 `json:"phaseCReactivePower"`
			VoltageAB                      float64 `json:"voltageAB"`
			VoltageBC                      float64 `json:"voltageBC"`
			VoltageCA                      float64 `json:"voltageCA"`
			PhaseAActiveEnergy             float64 `json:"phaseAActiveEnergy"`
			PhaseAReactiveEnergy           float64 `json:"phaseAReactiveEnergy"`
			PhaseBActiveEnergy             float64 `json:"phaseBActiveEnergy"`
			PhaseBReactiveEnergy           float64 `json:"phaseBReactiveEnergy"`
			PhaseCActiveEnergy             float64 `json:"phaseCActiveEnergy"`
			PhaseCReactiveEnergy           float64 `json:"phaseCReactiveEnergy"`
			EastronFlag                    int     `json:"eastronFlag"`
			EastronStatus                  int     `json:"eastronStatus"`
			ForwardTotalActiveEnergy       float64 `json:"forwardTotalActiveEnergy"`
			ReverseTotalActiveEnergy       float64 `json:"reverseTotalActiveEnergy"`
			TotalActivePower               float64 `json:"totalActivePower"`
			TotalApparentPower             float64 `json:"totalApparentPower"`
			PhaseAActivePower              float64 `json:"phaseAActivePower"`
			PhaseAApparentPower            float64 `json:"phaseAApparentPower"`
			PhaseAVoltage                  float64 `json:"phaseAVoltage"`
			PhaseACurrent                  float64 `json:"phaseACurrent"`
			PhaseAPowerFactor              float64 `json:"phaseAPowerFactor"`
			PhaseBActivePower              float64 `json:"phaseBActivePower"`
			PhaseBApparentPower            float64 `json:"phaseBApparentPower"`
			PhaseBVoltage                  float64 `json:"phaseBVoltage"`
			PhaseBCurrent                  float64 `json:"phaseBCurrent"`
			PhaseBPowerFactor              float64 `json:"phaseBPowerFactor"`
			PhaseCActivePower              float64 `json:"phaseCActivePower"`
			PhaseCApparentPower            float64 `json:"phaseCApparentPower"`
			PhaseCVoltage                  float64 `json:"phaseCVoltage"`
			PhaseCCurrent                  float64 `json:"phaseCCurrent"`
			PhaseCPowerFactor              float64 `json:"phaseCPowerFactor"`
		} `json:"datas"`
		Start    int  `json:"start"`
		HaveNext bool `json:"haveNext"`
	} `json:"obj"`
}

type GrowattNoahStatus struct {
	Result int         `json:"result"`
	Msg    interface{} `json:"msg"`
	Obj    struct {
		TotalBatteryPackSoc           string `json:"totalBatteryPackSoc"`
		Pac                           string `json:"pac"`
		WorkMode                      string `json:"workMode"`
		Ppv                           string `json:"ppv"`
		TotalBatteryPackChargingPower string `json:"totalBatteryPackChargingPower"`
		Status                        string `json:"status"`
	} `json:"obj"`
	Request interface{} `json:"request"`
}

type GrowattNoahTotals struct {
	Result int         `json:"result"`
	Msg    interface{} `json:"msg"`
	Obj    struct {
		MTotal    string `json:"mTotal"`
		MUnitText string `json:"mUnitText"`
		MToday    string `json:"mToday"`
		EacTotal  string `json:"eacTotal"`
		EacToday  string `json:"eacToday"`
	} `json:"obj"`
	Request interface{} `json:"request"`
}
