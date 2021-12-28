package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func SendTestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()
	return resp, string(respBody)
}

type want struct {
	responseStatusCode int
	responseParams     map[string]string
	responseBody       string
}

type request struct {
	url, method, body string
}

func TestGetPostNegative(t *testing.T) {
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "negative test #1. GET with empty url",
			request: request{
				url:    "/",
				method: http.MethodGet,
				body:   "",
			},
			want: want{
				responseStatusCode: http.StatusMethodNotAllowed,
				responseParams:     nil,
				responseBody:       "",
			},
		},
		{
			name: "negative test #2. GET with unresolved value",
			request: request{
				url:    "/RFGts",
				method: http.MethodGet,
				body:   "",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "Link not found",
			},
		},
		{
			name: "negative test #3. POST with empty body",
			request: request{
				url:    "/",
				method: http.MethodPost,
				body:   "",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "Request body is empty",
			},
		},
	}

	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, nil)
			assert.Equal(t, tt.want.responseStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.responseBody, TrimLastSymbols(body))
		})
	}
}

func TestShortenerHandlerPOSTMethod(t *testing.T) {

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "positive test #1. POST",
			request: request{
				url:    "/",
				method: http.MethodPost,
				body:   "http://yandex.ru",
			},
			want: want{
				responseStatusCode: http.StatusCreated,
				responseParams:     nil,
				responseBody:       "http://([a-zA-Z1-9]{5})",
			},
		},
	}
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, strings.NewReader(tt.request.body))

			assert.True(t, resp.StatusCode == tt.want.responseStatusCode)

			matched, _ := regexp.MatchString(tt.want.responseBody, TrimLastSymbols(body))
			assert.True(t, matched)
		})
	}
}

func TestShortenerHandlerGETMethodPositive(t *testing.T) {
	tests := []struct {
		name        string
		request     request
		originalURL string
		want        want
	}{
		{
			name:        "positive test #1. GET link",
			originalURL: "http://yandex.ru",
			request: request{
				url:    "/",
				method: http.MethodGet,
			},
			want: want{
				responseStatusCode: http.StatusTemporaryRedirect,
			},
		},
	}
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, shortLinkBody := SendTestRequest(t, ts, http.MethodPost, "/", strings.NewReader(tt.originalURL))
			shortLinksId := strings.Join(strings.Split(shortLinkBody, "/")[3:], "")

			getResp, _ := SendTestRequest(t, ts, tt.request.method, "/"+string(shortLinksId), nil)

			assert.True(t, getResp.Request.Response.Request.Response.StatusCode == tt.want.responseStatusCode)

			headers := getResp.Request.Response.Request.Response.Header.Get("Location")
			assert.Equal(t, headers, tt.originalURL)
		})
	}
}

func TrimLastSymbols(str string) string {
	return strings.TrimRight(string(str), "\n")
}
