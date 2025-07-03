package growatt_app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----- Test functions -----------------------------------------------------

type responseData struct {
	AnInt   int    `json:"an_int"`
	AString string `json:"a_string"`
}

func Test_postRequest(t *testing.T) {
	testData := []struct {
		url          string
		token        string
		data         url.Values
		responseBody any
		statusCode   int
		body         string
		errResult    bool
	}{
		{ // invalid body
			responseBody: responseData{},
			body:         `{invalid json}`,
			errResult:    true,
		},
		{ // error statusCode
			statusCode: 500,
			errResult:  true,
		},
		// { // io.ReadAll() fails - how...?
		// 	errResult: true,
		// },
		// { // http.Client.Do() fails - how...?
		// 	errResult: true,
		// },
		{ // request with token
			token:        "THETOKEN",
			responseBody: responseData{},
			body:         `{"an_int":33,"a_string":"another test"}`,
		},
		{ // http.NewRequest fails
			url:       "xyz://abc:8000f",
			token:     "",
			errResult: true,
		},
		{ // ok with body
			responseBody: responseData{},
			body:         `{"an_int":22,"a_string":"test"}`,
		},
		{ // ok without body
			body: "",
		},
	}

	for inx, td := range testData {
		t.Run(fmt.Sprintf("#%d", inx), func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
					tokenExists := len(td.token) > 0
					hdr, headerExists := req.Header["Authorization"]
					assert.Equal(t, tokenExists, headerExists)
					if headerExists {
						assert.Equal(t, fmt.Sprintf("Bearer %s", td.token), hdr[0])
					}

					ctHdr, ctHdrExists := req.Header["Content-Type"]
					assert.True(t, ctHdrExists)
					assert.Equal(t, "application/x-www-form-urlencoded", ctHdr[0])

					_, uaHdrExists := req.Header["User-Agent"]
					assert.True(t, uaHdrExists)

					wr.WriteHeader(td.statusCode)
					wr.Write([]byte(td.body))
				}))

			client := httpClient{client: server.Client()}
			if len(td.url) == 0 {
				td.url = server.URL
			}
			if td.statusCode == 0 {
				td.statusCode = 200
			}

			//			var data TokenResponse
			err := client.postForm(td.url, td.token, url.Values{}, td.responseBody)
			assert.Equal(t, td.errResult, err != nil)
			server.Close()
		})
	}
}
