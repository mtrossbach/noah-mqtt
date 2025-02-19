package growatt_web

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"noah-mqtt/internal/misc"
	"strings"
)

func (h *Client) postForm(url string, data url.Values, responseBody any) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: (HTTP %s) %s", resp.Status, string(b))
	}

	if responseBody != nil {
		if err := json.Unmarshal(b, &responseBody); err != nil {
			if strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
				if err := h.Login(); err != nil {
					slog.Error("could not re-login", slog.String("error", err.Error()))
					misc.Panic(err)
				}
				return h.postForm(url, data, responseBody)
			} else {
				return nil, err
			}
		}
	}

	return resp, nil
}
