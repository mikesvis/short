package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	"github.com/mikesvis/short/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, requestBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)

	client := ts.Client()

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testServer() *httptest.Server {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
	s := storage.NewStorage(c)
	h := NewHandler(c, s)
	l := logger.NewLogger()
	return httptest.NewServer(NewRouter(h, middleware.RequestResponseLogger(l)))
}

func TestShortRouter(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	startFull := "https://practicum.yandex.ru/"
	resp, startShort := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(startFull))
	resp.Body.Close()
	shortKey := string(startShort[len(startShort)-6:])

	type args struct {
		method string
		url    string
		body   io.Reader
	}

	type want struct {
		statusCode  int
		body        string
		contentType string
		location    string
	}

	tests := []struct {
		name string
		args args
		want
	}{
		{
			name: "Test POST / valid full url and store new short (201)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader("http://yandex.ru")},
			want: want{statusCode: http.StatusCreated, contentType: "text/plain"},
		}, {
			name: "Test POST / valid full url and get old short (409)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader(startFull)},
			want: want{statusCode: http.StatusConflict, body: startShort, contentType: "text/plain"},
		}, {
			name: "Test POST / invalid url on post (400)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader(":/ya")},
			want: want{statusCode: http.StatusBadRequest, body: "URL is not an URL format", contentType: "text/plain"},
		}, {
			name: "Test POST / invalid empty url (400)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader("")},
			want: want{statusCode: http.StatusBadRequest, body: "URL can not be empty", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten valid full url and store new short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"https://google.com"}`)},
			want: want{statusCode: http.StatusCreated, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten valid full url and get old short (409)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"` + startFull + `"}`)},
			want: want{statusCode: http.StatusConflict, body: startShort, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten invalid url on post (400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":":/ya"}`)},
			want: want{statusCode: http.StatusBadRequest, body: "URL is not an URL format", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten invalid empty url (400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":""}`)},
			want: want{statusCode: http.StatusBadRequest, body: "URL can not be empty", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten corrupted JSON(400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"}`)},
			want: want{statusCode: http.StatusBadRequest, body: "unexpected EOF", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten/batch valid full url and store new short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":"https://google.com"}]`)},
			want: want{statusCode: http.StatusCreated, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten/batch valid full url and get old short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":"` + startFull + `"}]`)},
			want: want{statusCode: http.StatusCreated, body: startShort, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten/batch invalid url on post (400)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":":/ya.ru"}]`)},
			want: want{statusCode: http.StatusBadRequest, body: "URL is not an URL format", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten/batch invalid empty url (400)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":""}]`)},
			want: want{statusCode: http.StatusBadRequest, body: "can not be empty", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten/batch corrupted JSON(400)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`{"url":"}`)},
			want: want{statusCode: http.StatusBadRequest, body: "unexpected EOF", contentType: "text/plain"},
		}, {
			name: "Test GET / success get (307 -> redirect -> 200)",
			args: args{method: http.MethodGet, url: shortKey},
			want: want{statusCode: http.StatusOK},
		}, {
			name: "Test GET / fail (400)",
			args: args{method: http.MethodGet, url: "/iddQd-doom-slayer"},
			want: want{statusCode: http.StatusBadRequest, body: "full url is not found", contentType: "text/plain"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.args.method, tt.args.url, tt.args.body)
			resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Contains(t, get, tt.want.body)
			if len(tt.want.contentType) > 0 {
				assert.Contains(t, resp.Header.Get("Content-type"), tt.want.contentType)
			}

			if len(tt.want.location) > 0 {
				assert.Contains(t, resp.Header.Get("Location"), tt.want.location)
			}
		})
	}
}
