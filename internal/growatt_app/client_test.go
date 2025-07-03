package growatt_app

import (
	"errors"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----- Mocks --------------------------------------------------------------

// MockHttpClient implements HttpClient
type MockHttpClient struct {
	mock.Mock
}

func (h *MockHttpClient) postForm(url string, token string, data url.Values, responseBody any) error {
	args := h.Called(url, token, data, responseBody)
	return args.Error(0)
}

// ----- Test functions -----------------------------------------------------

func setupMocks(t *testing.T) (*MockHttpClient, *Client) {
	mockHttpClient := MockHttpClient{}
	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

	client := Client{
		client:    &mockHttpClient,
		serverUrl: "https://server-api.growatt.com",
		username:  "user",
		password:  "secret",
		jar:       jar,
	}

	return &mockHttpClient, &client
}

func Test_postForm_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_AnyError(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("some error"))

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLogin(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value"))

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*TokenResponse)
		responseBody.Token = "THE_TOKEN"
	}).Return(nil)

	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == "user" &&
			m.Get("password") == "secret" &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		"THE_TOKEN",
		expectedMatch,
		&LoginResult{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*LoginResult)
		responseBody.Back.Success = true
	}).Return(nil)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"THE_TOKEN",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLoginFails(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value"))

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*TokenResponse)
		responseBody.Token = "THE_TOKEN"
	}).Return(nil)

	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == "user" &&
			m.Get("password") == "secret" &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		"THE_TOKEN",
		expectedMatch,
		&LoginResult{},
	).Return(errors.New("login failed"))

	defer func() {
		if r := recover(); r != nil {
			mockHttpClient.AssertExpectations(t)
		}
	}()

	client.postForm("http://someurl", url.Values{}, nil)
	t.Errorf("Test failed, panic was expected")
}

func TestLogin_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*TokenResponse)
		responseBody.Token = "THE_TOKEN"
	}).Return(nil)

	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == "user" &&
			m.Get("password") == "secret" &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		"THE_TOKEN",
		expectedMatch,
		&LoginResult{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*LoginResult)
		responseBody.Back.Success = true
	}).Return(nil)

	err := client.Login()
	assert.NoError(t, err)
	assert.Equal(t, "THE_TOKEN", client.token)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_NoSuccess(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*TokenResponse)
		responseBody.Token = "THE_TOKEN"
	}).Return(nil)

	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == "user" &&
			m.Get("password") == "secret" &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		"THE_TOKEN",
		expectedMatch,
		&LoginResult{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*LoginResult)
		responseBody.Back.Success = false
	}).Return(nil)

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_loginGetTokenFail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Return(errors.New("login error"))

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_newTwoLoginFail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {"SHINEuser"},
			"password": {"secret"},
			"lan":      {"1"},
		},
		&TokenResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*TokenResponse)
		responseBody.Token = "THE_TOKEN"
	}).Return(nil)

	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == "user" &&
			m.Get("password") == "secret" &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		"THE_TOKEN",
		expectedMatch,
		&LoginResult{},
	).Return(errors.New("newTwoLoginAPIV2 fail"))

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestGetPlantList_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoPlantAPI.do?op=getAllPlantListTwo",
		"",
		url.Values{
			"plantStatus": {""},
			"pageSize":    {"20"},
			"language":    {"1"},
			"toPageNum":   {"1"},
			"order":       {"1"},
		},
		&PlantListV2{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*PlantListV2)
		*responseBody = PlantListV2{
			PlantList: []struct {
				ID int `json:"id"`
			}{
				{ID: 1},
				{ID: 2},
			},
		}
	}).Return(nil)

	data, err := client.GetPlantList()

	assert.NoError(t, err)
	assert.Equal(t, PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: 1},
			{ID: 2},
		},
	}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetPlantList_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/newTwoPlantAPI.do?op=getAllPlantListTwo",
		"",
		url.Values{
			"plantStatus": {""},
			"pageSize":    {"20"},
			"language":    {"1"},
			"toPageNum":   {"1"},
			"order":       {"1"},
		},
		&PlantListV2{},
	).Return(errors.New("newTwoPlantAPI fail"))

	data, err := client.GetPlantList()
	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/noah/isPlantNoahSystem",
		"",
		url.Values{
			"plantId": {"2"},
		},
		&NoahPlantInfo{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*NoahPlantInfo)
		*responseBody = NoahPlantInfo{
			ResponseContainerV2: ResponseContainerV2[NoahPlantInfoObj]{
				Msg:    "",
				Result: 0,
				Obj: NoahPlantInfoObj{
					IsPlantHaveNexa: true,
				},
			},
		}
	}).Return(nil)

	data, err := client.GetNoahPlantInfo("2")

	assert.NoError(t, err)
	assert.Equal(t, NoahPlantInfo{
		ResponseContainerV2: ResponseContainerV2[NoahPlantInfoObj]{
			Msg:    "",
			Result: 0,
			Obj: NoahPlantInfoObj{
				IsPlantHaveNexa: true,
			},
		},
	}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_NoNexa(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/noah/isPlantNoahSystem",
		"",
		url.Values{
			"plantId": {"2"},
		},
		&NoahPlantInfo{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*NoahPlantInfo)
		*responseBody = NoahPlantInfo{
			ResponseContainerV2: ResponseContainerV2[NoahPlantInfoObj]{
				Msg:    "",
				Result: 0,
				Obj: NoahPlantInfoObj{
					IsPlantHaveNexa: false,
				},
			},
		}
	}).Return(nil)

	data, err := client.GetNoahPlantInfo("2")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/noah/isPlantNoahSystem",
		"",
		url.Values{
			"plantId": {"2"},
		},
		&NoahPlantInfo{},
	).Return(errors.New("isPlantNoahSystem fail"))

	data, err := client.GetNoahPlantInfo("2")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahStatus_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getSystemStatus",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&NoahStatus{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*NoahStatus)
		*responseBody = NoahStatus{}
	}).Return(nil)

	data, err := client.GetNoahStatus("serial123")

	assert.NoError(t, err)
	assert.Equal(t, NoahStatus{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahStatus_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getSystemStatus",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&NoahStatus{},
	).Return(errors.New("getSystemStatus fail"))

	data, err := client.GetNoahStatus("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahInfo_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getNexaInfoBySn",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&NexaInfo{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*NexaInfo)
		*responseBody = NexaInfo{}
	}).Return(nil)

	data, err := client.GetNoahInfo("serial123")

	assert.NoError(t, err)
	assert.Equal(t, NexaInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahInfo_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getNexaInfoBySn",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&NexaInfo{},
	).Return(errors.New("getSystemStatus fail"))

	data, err := client.GetNoahInfo("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetBatteryData_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getBatteryData",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&BatteryInfo{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*BatteryInfo)
		*responseBody = BatteryInfo{}
	}).Return(nil)

	data, err := client.GetBatteryData("serial123")

	assert.NoError(t, err)
	assert.Equal(t, BatteryInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetBatteryData_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getBatteryData",
		"",
		url.Values{
			"deviceSn": {"serial123"},
		},
		&BatteryInfo{},
	).Return(errors.New("getBatteryData fail"))

	data, err := client.GetBatteryData("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestSetSystemOutputPower_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"system_out_put_power"},
			"param1":    {"0"},
			"param2":    {"200"},
		},
		&SetResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*SetResponse)
		*responseBody = SetResponse{}
	}).Return(nil)

	err := client.SetSystemOutputPower("serial123", 0, 200)

	assert.NoError(t, err)
	//	assert.Equal(t, BatteryInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestSetSystemOutputPower_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"system_out_put_power"},
			"param1":    {"0"},
			"param2":    {"200"},
		},
		&SetResponse{},
	).Return(errors.New("noahDeviceApi/nexa/set system_out_put_power fail"))

	err := client.SetSystemOutputPower("serial123", 0, 200)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingSoc_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"charging_soc"},
			"param1":    {"85"},
			"param2":    {"15"},
		},
		&SetResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*SetResponse)
		*responseBody = SetResponse{}
	}).Return(nil)

	err := client.SetChargingSoc("serial123", 85, 15)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingSoc_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"charging_soc"},
			"param1":    {"85"},
			"param2":    {"15"},
		},
		&SetResponse{},
	).Return(errors.New("noahDeviceApi/nexa/set charging_soc fail"))

	err := client.SetChargingSoc("serial123", 85, 15)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"allow_grid_charging"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*SetResponse)
		*responseBody = SetResponse{}
	}).Return(nil)

	err := client.SetAllowGridCharging("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"allow_grid_charging"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Return(errors.New("noahDeviceApi/nexa/set allow_grid_charging fail"))

	err := client.SetAllowGridCharging("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"grid_connection_control"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*SetResponse)
		*responseBody = SetResponse{}
	}).Return(nil)

	err := client.SetGridConnectionControl("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"grid_connection_control"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Return(errors.New("noahDeviceApi/nexa/set grid_connection_control fail"))

	err := client.SetGridConnectionControl("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetACCouplePowerControl_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"ac_couple_power_control"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Run(func(args mock.Arguments) {
		responseBody := args.Get(3).(*SetResponse)
		*responseBody = SetResponse{}
	}).Return(nil)

	err := client.SetACCouplePowerControl("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetACCouplePowerControl_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {"serial123"},
			"type":      {"ac_couple_power_control"},
			"param1":    {"1"},
		},
		&SetResponse{},
	).Return(errors.New("noahDeviceApi/nexa/set ac_couple_power_control fail"))

	err := client.SetACCouplePowerControl("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}
