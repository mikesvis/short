package server

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/jwt"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	"github.com/mikesvis/short/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, requestBody io.Reader, cookies []*http.Cookie) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	if len(cookies) > 0 {
		client.Jar.SetCookies(req.URL, cookies)
	}

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
	l, _ := logger.NewLogger()
	s, _ := storage.NewStorage(c, l)
	h := NewHandler(c, s)
	return httptest.NewServer(NewRouter(h, middleware.RequestResponseLogger(l)))
}

// я немного очумел пока это сделал
func generateTestCookiesByUser(userID string) []*http.Cookie {
	startTokenString, _ := jwt.CreateTokenString(userID, time.Now().Add(5*time.Minute))
	startCookie := middleware.CreateAuthCookie(startTokenString, time.Now().Add(5*time.Minute))
	cookies := []*http.Cookie{}
	cookies = append(cookies, startCookie)
	return cookies
}

func TestShortRouter(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	startFull := "https://practicum.yandex.ru/"
	resp, startShort := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(startFull), generateTestCookiesByUser("DoomGuy"))
	resp.Body.Close()
	shortKey := string(startShort[len(startShort)-6:])

	type args struct {
		method  string
		url     string
		body    io.Reader
		cookies []*http.Cookie
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
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader("http://yandex.ru"), cookies: nil},
			want: want{statusCode: http.StatusCreated, contentType: "text/plain"},
		}, {
			name: "Test POST / valid full url and get old short (409)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader(startFull), cookies: nil},
			want: want{statusCode: http.StatusConflict, body: startShort, contentType: "text/plain"},
		}, {
			name: "Test POST / invalid url on post (400)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader(":/ya"), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "URL is not an URL format", contentType: "text/plain"},
		}, {
			name: "Test POST / invalid empty url (400)",
			args: args{method: http.MethodPost, url: "/", body: strings.NewReader(""), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "URL can not be empty", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten valid full url and store new short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"https://google.com"}`), cookies: nil},
			want: want{statusCode: http.StatusCreated, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten valid full url and get old short (409)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"` + startFull + `"}`), cookies: nil},
			want: want{statusCode: http.StatusConflict, body: startShort, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten invalid url on post (400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":":/ya"}`), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "URL is not an URL format", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten invalid empty url (400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":""}`), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "URL can not be empty", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten corrupted JSON(400)",
			args: args{method: http.MethodPost, url: "/api/shorten", body: strings.NewReader(`{"url":"}`), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "unexpected EOF", contentType: "text/plain"},
		}, {
			name: "Test POST /api/shorten/batch valid full url and store new short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":"https://google.com"}]`), cookies: nil},
			want: want{statusCode: http.StatusCreated, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten/batch valid full url and get old short (201)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`[{"correlation_id":"1","original_url":"` + startFull + `"}]`), cookies: nil},
			want: want{statusCode: http.StatusCreated, body: startShort, contentType: "application/json"},
		}, {
			name: "Test POST /api/shorten/batch corrupted JSON(400)",
			args: args{method: http.MethodPost, url: "/api/shorten/batch", body: strings.NewReader(`{"url":"}`), cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "unexpected EOF", contentType: "text/plain"},
		}, {
			name: "Test GET / success get (307 -> redirect -> 200)",
			args: args{method: http.MethodGet, url: shortKey, cookies: nil},
			want: want{statusCode: http.StatusOK},
		}, {
			name: "Test GET / fail (400)",
			args: args{method: http.MethodGet, url: "/iddQd-doom-slayer", cookies: nil},
			want: want{statusCode: http.StatusBadRequest, body: "full url is not found", contentType: "text/plain"},
		}, {
			name: "Test GET /api/user/urls with list (200)",
			args: args{method: http.MethodGet, url: "/api/user/urls", cookies: generateTestCookiesByUser("DoomGuy")},
			want: want{statusCode: http.StatusOK, contentType: "application/json"},
		}, {
			name: "Test GET /api/user/urls with empty list(204)",
			args: args{method: http.MethodGet, url: "/api/user/urls", cookies: generateTestCookiesByUser("Heretic")},
			want: want{statusCode: http.StatusNoContent},
		}, {
			name: "Test GET /api/user/urls with unauthozied(401)",
			args: args{method: http.MethodGet, url: "/api/user/urls", cookies: nil},
			want: want{statusCode: http.StatusUnauthorized},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.args.method, tt.args.url, tt.args.body, tt.args.cookies)
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
