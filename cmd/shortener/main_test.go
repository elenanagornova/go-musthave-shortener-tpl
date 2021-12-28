package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func SendTestRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

//
func TestRouter(t *testing.T) {
	type want struct {
		responseStatusCode int
		responseParams     map[string]string
		responseBody       string
	}

	type request struct {
		url, method, body string
	}

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
	}

	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, body := SendTestRequest(t, ts, tt.request.method, tt.request.url)
			assert.Equal(t, tt.want.responseStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.responseBody, strings.TrimRight(string(body), "\n"))
		})
	}
}
func TestShortenerHandlerPOSTMethod(t *testing.T) {
	type want struct {
		responseStatusCode int
		responseParams     map[string]string
		responseBody       string
	}

	type request struct {
		url, method, body string
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "negative test #1. POST with empty body",
			request: request{
				url:    "/",
				method: http.MethodPost,
				body:   "",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "",
			},
		},
		{
			name: "positive test #2. POST",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(MakeShortLink)
			// запуск сервера
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.True(t, res.StatusCode == tt.want.responseStatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			resBodyAsStr := strings.TrimRight(string(resBody), "\n")
			// если запрос POST и в ответе не пусто и запрос был принят, тогда проверяем ответ по regex
			if tt.want.responseBody != "" && tt.request.method == http.MethodPost && res.StatusCode == http.StatusCreated {
				matched, _ := regexp.MatchString(tt.want.responseBody, resBodyAsStr)
				assert.True(t, matched)
			} else {
				assert.True(t, strings.HasPrefix(resBodyAsStr, tt.want.responseBody))
			}
		})
	}
}

func TestShortenerHandlerGETMethod(t *testing.T) {
	type want struct {
		responseStatusCode int
		responseParams     map[string]string
		responseBody       string
	}

	type request struct {
		url, method, body string
	}

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
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "The path is missing",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetLinkByID)
			// запуск сервера
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.True(t, res.StatusCode == tt.want.responseStatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			resBodyAsStr := strings.TrimRight(string(resBody), "\n")
			// если запрос POST и в ответе не пусто и запрос был принят, тогда проверяем ответ по regex
			if tt.want.responseBody != "" && tt.request.method == http.MethodPost && res.StatusCode == http.StatusCreated {
				matched, _ := regexp.MatchString(tt.want.responseBody, resBodyAsStr)
				assert.True(t, matched)
			} else {
				assert.True(t, strings.HasPrefix(resBodyAsStr, tt.want.responseBody))
			}
		})
	}
}

func TestShortenerHandlerGETMethodPositive(t *testing.T) {
	type want struct {
		responseStatusCode int
		responseParams     map[string]string
	}

	type request struct {
		url, method string
	}

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortLinks := sendPostRequestForTests(tt.originalURL)

			//отправляем гет запрос на данный реквест
			getLinkrequest := httptest.NewRequest(tt.request.method, shortLinks, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetLinkByID)
			// запуск сервера
			h.ServeHTTP(w, getLinkrequest)
			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.True(t, res.StatusCode == tt.want.responseStatusCode)

			headers := res.Header.Get("Location")
			assert.Equal(t, headers, tt.originalURL)
		})
	}
}

func sendPostRequestForTests(originalURL string) string {
	createLinkRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(MakeShortLink)
	// запуск сервера
	h.ServeHTTP(w, createLinkRequest)
	shortLink := w.Result()
	defer shortLink.Body.Close()
	shortLinkBody, err := io.ReadAll(shortLink.Body)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(string(shortLinkBody), "\n")
}
