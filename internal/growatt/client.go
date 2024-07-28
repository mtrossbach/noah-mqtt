package growatt

import (
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/backoff/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	client       *http.Client
	username     string
	password     string
	userAgent    string
	userId       string
	jar          *cookiejar.Jar
	loginBackoff backoff.Policy
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
		loginBackoff: backoff.Exponential(
			backoff.WithMinInterval(10*time.Second),
			backoff.WithMaxInterval(10*time.Minute),
			backoff.WithJitterFactor(0.05),
		),
	}
}

func (h *Client) Login() error {
	resp, err := h.postForm("https://openapi.growatt.com/newTwoLoginAPI.do", url.Values{
		"userName": {h.username},
		"password": {h.password},
	})
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data LoginResult
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if resp.StatusCode != 200 || !data.Back.Success {
		slog.Error("login failed", slog.String("data", string(b)))
		slog.Info("waiting before exiting")
		<-time.After(60 * time.Second)
		panic("login failed")
	}

	h.userId = fmt.Sprintf("%d", data.Back.User.ID)
	return nil
}

func (h *Client) GetPlantList() (*PlantList, error) {
	resp, err := h.client.Get(fmt.Sprintf("https://openapi.growatt.com/PlantListAPI.do?userId=%s", h.userId))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get plant list failed: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PlantList
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	if !result.Back.Success {
		return nil, fmt.Errorf("get plant list failed")
	}

	return &result, nil
}

func (h *Client) GetNoahPlantInfo(plantId string) (*NoahPlantInfo, error) {
	resp, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/isPlantNoahSystem", url.Values{
		"plantId": {plantId},
	})

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get noah serial failed: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result NoahPlantInfo
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *Client) GetNoahStatus(serialNumber string) (*NoahStatus, error) {
	resp, err := h.postForm("https://openapi.growatt.com/noahDeviceApi/noah/getSystemStatus", url.Values{
		"deviceSn": {serialNumber},
	})

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetch noah status failed: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data NoahStatus
	if err := json.Unmarshal(b, &data); err != nil {
		if strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
			if err := h.Login(); err != nil {
				return nil, err
			}
			return h.GetNoahStatus(serialNumber)
		} else {
			slog.Error("could not parse json", slog.String("error", err.Error()), slog.String("data", string(b)))
			return nil, err
		}
	}

	return &data, nil
}

func (h *Client) postForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return h.client.Do(req)
}
