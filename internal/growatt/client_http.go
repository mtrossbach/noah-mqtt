package growatt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (h *Client) postForm(url string, data url.Values, responseBody any) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(h.token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	}
	req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 9; Mi A1 MIUI/V10.0.24.0.PDHMIXM)")
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
					<-time.After(60 * time.Second)
					panic(err)
				}
				return h.postForm(url, data, responseBody)
			} else {
				return nil, err
			}
		}
	}

	return resp, nil
}
