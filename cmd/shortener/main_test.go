package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-musthave-shortener-tpl/internal/shortener"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func SendTestRequest(t *testing.T, ts *httptest.Server, method, path, contentType string, body io.Reader) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()
	return resp, string(respBody)
}

type want struct {
	responseStatusCode  int
	responseParams      map[string]string
	responseContentType string
	responseBody        string
}

type request struct {
	url, method, contentType, body string
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
				responseStatusCode: http.StatusNotFound,
				responseParams:     nil,
				responseBody:       "link not found",
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

	service := shortener.New(addr)

	r := NewRouter(service)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, "", nil)
			defer resp.Body.Close()
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
	service := shortener.New(addr)
	r := NewRouter(service)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, "", strings.NewReader(tt.request.body))
			defer resp.Body.Close()

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
	service := shortener.New("http://localhost:8080")
	r := NewRouter(service)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postResp, shortLinkBody := SendTestRequest(t, ts, http.MethodPost, "/", "", strings.NewReader(tt.originalURL))
			defer postResp.Body.Close()

			shortLinksID := strings.Join(strings.Split(shortLinkBody, "/")[3:], "")

			getResp, _ := SendTestRequest(t, ts, tt.request.method, "/"+string(shortLinksID), "", nil)
			defer getResp.Body.Close()

			assert.True(t, getResp.StatusCode == tt.want.responseStatusCode)

			headers := getResp.Header.Get("Location")
			assert.Equal(t, headers, tt.originalURL)
		})
	}
}

func TestMakeShortenLinkPOSTMethodPositive(t *testing.T) {

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "positive test #1. POST /api/shorten",
			request: request{
				url:         "/api/shorten",
				contentType: "application/json",
				method:      http.MethodPost,
				body:        `{"url": "http://yandex.ru"}`,
			},
			want: want{
				responseStatusCode:  201,
				responseParams:      nil,
				responseContentType: "application/json",
				responseBody:        "http://([a-zA-Z1-9]{5})",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := shortener.New(addr)
			r := NewRouter(service)
			ts := httptest.NewServer(r)
			defer ts.Close()
			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, tt.request.contentType, strings.NewReader(tt.request.body))
			defer resp.Body.Close()

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.responseContentType)
			assert.Equal(t, tt.want.responseStatusCode, resp.StatusCode)

			matched, _ := regexp.MatchString(tt.want.responseBody, TrimLastSymbols(body))
			assert.True(t, matched)
		})
	}
}
func TestMakeShortenLinkPOSTMethodNegative(t *testing.T) {

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "negative test #1. POST /api/shorten without contentType",
			request: request{
				url:         "/api/shorten",
				contentType: "",
				method:      http.MethodPost,
				body:        `{"url": "http://yandex.ru"}`,
			},
			want: want{
				responseStatusCode: http.StatusUnsupportedMediaType,
				responseParams:     nil,
				responseBody:       "Content Type is not application/json",
			},
		},
		{
			name: "negative test #2. POST /api/shorten with wrong contentType",
			request: request{
				url:         "/api/shorten",
				contentType: "text/plain",
				method:      http.MethodPost,
				body:        `{"url": "http://yandex.ru"}`,
			},
			want: want{
				responseStatusCode: http.StatusUnsupportedMediaType,
				responseParams:     nil,
				responseBody:       "Content Type is not application/json",
			},
		},
		{
			name: "negative test #3. POST /api/shorten with wrong body",
			request: request{
				url:         "/api/shorten",
				contentType: "application/json",
				method:      http.MethodPost,
				body:        `{"link": }`,
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "Something wrong with request",
			},
		},
		{
			name: "negative test #4. POST /api/shorten with empty body",
			request: request{
				url:         "/api/shorten",
				contentType: "application/json",
				method:      http.MethodPost,
				body:        "",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "Something wrong with request",
			},
		},
	}
	service := shortener.New(addr)
	r := NewRouter(service)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url, tt.request.contentType, strings.NewReader(tt.request.body))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.responseStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.responseBody, TrimLastSymbols(body))

		})
	}
}
func TrimLastSymbols(str string) string {
	return strings.TrimRight(str, "\n")
}
