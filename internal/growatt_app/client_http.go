package growatt_app

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type httpClient struct {
	client *http.Client
}

type HttpClient interface {
	postForm(url string, token string, data url.Values, responseBody any) error
}

func (h *httpClient) postForm(url string, token string, data url.Values, responseBody any) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("http.NewRequest failed (app)", slog.String("error", err.Error()))
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/27.0 Chrome/125.0.0.0 Mobile Safari/537.36")
	resp, err := h.client.Do(req)
	if err != nil {
		slog.Error("http.Client.Do failed (app)", slog.String("error", err.Error()))
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("io.ReadAll failed (app)", slog.String("error", err.Error()))
		return err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("request failed: (HTTP %s) %s", resp.Status, string(b))
		slog.Error("StatusCode != 200 (app)", slog.String("error", err.Error()))
		return err
	}

	if responseBody != nil {
		if err := json.Unmarshal(b, &responseBody); err != nil {
			slog.Error("json.Unmarshal failed (app)", slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}
