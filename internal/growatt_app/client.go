package growatt_app

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"nexa-mqtt/internal/misc"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	client    HttpClient
	serverUrl string
	username  string
	password  string
	userAgent string
	userId    string
	token     string
	jar       *cookiejar.Jar
}

func newClient(serverUrl string, username string, password string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("could not create cookie jar", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	slog.Info("setting server url (app)", slog.String("url", serverUrl))

	return &Client{
		client: &httpClient{
			client: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           jar,
				Timeout:       10 * time.Second,
			},
		},
		serverUrl: serverUrl,
		username:  username,
		password:  hashPassword(password),
		jar:       jar,
	}
}

func (h *Client) postForm(url string, data url.Values, responseBody any) error {
	err := h.client.postForm(url, h.token, data, responseBody)
	if err != nil {
		notLoggedIn := strings.Contains(err.Error(), "Dear user, you have not login to the system") ||
			strings.Contains(err.Error(), "invalid character '<' looking for beginning of value")

		if notLoggedIn {
			slog.Warn("re-login", slog.String("error", err.Error()))
			if err := h.Login(); err != nil {
				slog.Error("could not re-login", slog.String("error", err.Error()))
				misc.Panic(err)
			}
			return h.postForm(url, data, responseBody)
		} else {
			return err
		}
	}

	return nil
}

func (h *Client) loginGetToken() error {
	var data TokenResponse
	if err := h.postForm("https://evcharge.growatt.com/ocpp/user", url.Values{
		"cmd":      {"shineLogin"},
		"userId":   {fmt.Sprintf("SHINE%s", h.username)},
		"password": {h.password},
		"lan":      {"1"},
	}, &data); err != nil {
		return err
	}

	h.token = data.Token
	return nil
}

func (h *Client) Login() error {
	if err := h.loginGetToken(); err != nil {
		return err
	}

	var data LoginResult
	if err := h.postForm(h.serverUrl+"/newTwoLoginAPIV2.do", url.Values{
		"userName":          {h.username},
		"password":          {h.password},
		"newLogin":          {"1"},
		"phoneType":         {"android"},
		"shinephoneVersion": {"8.3.0.2"},
		"phoneSn":           {uuid.New().String()},
		"ipvcpc":            {ipvcpc(h.username)},
		"language":          {"1"},
		"systemVersion":     {"15"},
		"phoneModel":        {"Mi A1"},
		"loginTime":         {time.Now().Format(time.DateTime)},
		"appType":           {"ShinePhone"},
		"timestamp":         {timestamp()},
	}, &data); err != nil {
		return err
	}

	if !data.Back.Success {
		return fmt.Errorf("login failed: %s", data.Back.Msg)
	}

	h.userId = fmt.Sprintf("%d", data.Back.User.ID)
	return nil
}

func (h *Client) GetPlantList() (*PlantListV2, error) {
	var data PlantListV2
	if err := h.postForm(h.serverUrl+"/newTwoPlantAPI.do?op=getAllPlantListTwo", url.Values{
		"plantStatus": {""},
		"pageSize":    {"20"},
		"language":    {"1"},
		"toPageNum":   {"1"},
		"order":       {"1"},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahPlantInfo(plantId string) (*NoahPlantInfo, error) {
	var data NoahPlantInfo
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/noah/isPlantNoahSystem", url.Values{
		"plantId": {plantId},
	}, &data); err != nil {
		return nil, err
	}

	if !data.Obj.IsPlantHaveNexa {
		return nil, errors.New("No NEXA device")
	}

	return &data, nil
}

func (h *Client) GetNoahStatus(serialNumber string) (*NoahStatus, error) {
	var data NoahStatus
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getSystemStatus", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahInfo(serialNumber string) (*NexaInfo, error) {
	var data NexaInfo
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getNexaInfoBySn", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) GetBatteryData(serialNumber string) (*BatteryInfo, error) {
	var data BatteryInfo
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getBatteryData", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) SetSystemOutputPower(serialNumber string, mode int, power float64) error {
	p := math.Max(0, math.Min(800, power))
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"system_out_put_power"},
		"param1":    {fmt.Sprintf("%d", mode)},
		"param2":    {fmt.Sprintf("%.0f", p)},
	}, &data); err != nil {
		return err
	}

	return nil
}

func (h *Client) SetChargingSoc(serialNumber string, chargingLimit float64, dischargeLimit float64) error {
	c := math.Max(70, math.Min(100, chargingLimit))
	d := math.Max(0, math.Min(30, dischargeLimit))
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"charging_soc"},
		"param1":    {fmt.Sprintf("%.0f", c)},
		"param2":    {fmt.Sprintf("%.0f", d)},
	}, &data); err != nil {
		return err
	}

	return nil
}

func (h *Client) SetAllowGridCharging(serialNumber string, allow int) error {
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"allow_grid_charging"},
		"param1":    {fmt.Sprintf("%d", allow)},
	}, &data); err != nil {
		return err
	}

	return nil
}

func (h *Client) SetGridConnectionControl(serialNumber string, offlineEnable int) error {
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"grid_connection_control"},
		"param1":    {fmt.Sprintf("%d", offlineEnable)},
	}, &data); err != nil {
		return err
	}

	return nil
}

func (h *Client) SetACCouplePowerControl(serialNumber string, _1000WEnable int) error {
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"ac_couple_power_control"},
		"param1":    {fmt.Sprintf("%d", _1000WEnable)},
	}, &data); err != nil {
		return err
	}

	return nil
}
