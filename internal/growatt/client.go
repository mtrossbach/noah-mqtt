package growatt

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	client    *http.Client
	username  string
	password  string
	userAgent string
	userId    string
	jar       *cookiejar.Jar
}

func NewClient(username string, password string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &Client{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           jar,
			Timeout:       30 * time.Second,
		},
		username: username,
		password: hashPassword(password),
		jar:      jar,
	}
}

func (h *Client) Login() error {
	var data LoginResult
	if _, err := h.postForm("https://openapi.growatt.com/newTwoLoginAPI.do", url.Values{
		"userName": {h.username},
		"password": {h.password},
	}, &data); err != nil {
		return err
	}

	if !data.Back.Success {
		return fmt.Errorf("login failed: %s", data.Back.Msg)
	}

	h.userId = fmt.Sprintf("%d", data.Back.User.ID)
	return nil
}

func (h *Client) GetPlantList() (*PlantList, error) {
	var data PlantList
	if _, err := h.get("https://openapi.growatt.com/PlantListAPI.do", url.Values{
		"userId": {h.userId},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahPlantInfo(plantId string) (*NoahPlantInfo, error) {
	var data NoahPlantInfo
	if _, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/isPlantNoahSystem", url.Values{
		"plantId": {plantId},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahStatus(serialNumber string) (*NoahStatus, error) {
	var data NoahStatus
	if _, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/getSystemStatus", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahInfo(serialNumber string) (*NoahInfo, error) {
	var data NoahInfo
	if _, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/getNoahInfoBySn", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) GetBatteryData(serialNumber string) (*BatteryInfo, error) {
	var data BatteryInfo
	if _, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/getBatteryData", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
