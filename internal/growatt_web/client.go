package growatt_web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"noah-mqtt/internal/misc"
	"time"
)

type Client struct {
	client    *http.Client
	serverUrl string
	username  string
	password  string
	jar       *cookiejar.Jar
}

func newClient(serverUrl string, username string, password string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("could not create cookie jar", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	slog.Info("setting server url (web)", slog.String("url", serverUrl))

	return &Client{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           jar,
			Timeout:       10 * time.Second,
		},
		serverUrl: serverUrl,
		username:  username,
		password:  password,
		jar:       jar,
	}
}

func (h *Client) Login() error {
	var result GrowattResult
	if _, err := h.postForm(h.serverUrl+"/login", url.Values{
		"account":  {h.username},
		"password": {h.password},
	}, &result); err != nil {
		return err
	}
	return nil
}

func (h *Client) GetPlantList() ([]GrowattPlant, error) {
	var result []GrowattPlant
	if _, err := h.postForm(h.serverUrl+"/index/getPlantListTitle", url.Values{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Client) GetPlantDevices(plantId string) (*GrowattPlantDevices, error) {
	var result GrowattPlantDevices
	if _, err := h.postForm(h.serverUrl+"/panel/getDevicesByPlantList", url.Values{
		"plantId":  {plantId},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahList(plantId int) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if _, err := h.postForm(h.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahDetails(plantId int, serial string) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if _, err := h.postForm(h.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"deviceSn": {serial},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahHistory(serial string, startDate string, endDate string) (*GrowattNoahHistory, error) {
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}
	var result GrowattNoahHistory
	if _, err := h.postForm(h.serverUrl+"/device/getNoahHistory", url.Values{
		"deviceSn":  {serial},
		"start":     {"0"},
		"startDate": {startDate},
		"endDate":   {endDate},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahStatus(plantId int, serial string) (*GrowattNoahStatus, error) {
	var result GrowattNoahStatus
	if _, err := h.postForm(fmt.Sprintf(h.serverUrl+"/panel/noah/getNoahStatusData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahTotals(plantId int, serial string) (*GrowattNoahTotals, error) {
	var result GrowattNoahTotals
	if _, err := h.postForm(fmt.Sprintf(h.serverUrl+"/panel/noah/getNoahTotalData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
