package server

import (
	_context "context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	return &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
}

func TestGetFullURL(t *testing.T) {
	c := testConfig()
	s := inmemory.NewInMemory()
	s.Store(_context.Background(), domain.URL{
		Full:  "http://www.yandex.ru/verylongpath",
		Short: "short",
	})

	type want struct {
		statusCode  int
		newLocation string
		wantError   bool
		body        string
	}
	type request struct {
		methhod string
		target  string
	}

	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Find full url (307)",
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				newLocation: "http://www.yandex.ru/verylongpath",
				wantError:   false,
			},
			request: request{
				methhod: "GET",
				target:  "/short",
			},
		}, {
			name: "Full url does not exist (400)",
			want: want{
				statusCode:  http.StatusBadRequest,
				newLocation: "",
				wantError:   true,
				body:        "full url is not found",
			},
			request: request{
				methhod: "GET",
				target:  "http://example.com/short",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.methhod, tt.request.target, nil)
			w := httptest.NewRecorder()
			handler := NewHandler(c, s)
			handle := http.HandlerFunc(handler.GetFullURL)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if !tt.want.wantError {
				assert.Contains(t, result.Header.Values("Location"), tt.want.newLocation)
				return
			}

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			require.Contains(t, string(response), tt.want.body)
		})
	}
}

func TestCreateShortURLText(t *testing.T) {
	c := testConfig()
	s := inmemory.NewInMemory()
	ctx := _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy")
	s.Store(ctx, domain.URL{
		Full:   "http://www.yandex.ru/verylongpath",
		UserID: "Doomguy",
		Short:  "short",
	})

	type want struct {
		contentType string
		statusCode  int
		isNew       bool
		wantError   bool
		body        string
	}
	type request struct {
		method string
		target string
		body   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Get old short url from full (409)",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
				isNew:       false,
				wantError:   false,
				body:        string(c.BaseURL) + "/short",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/verylongpath",
			},
		}, {
			name: "Create new short url from full (201)",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				isNew:       true,
				wantError:   false,
				body:        string(c.BaseURL),
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/very",
			},
		}, {
			name: "Empty body (400)",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				wantError:   true,
				body:        "URL can not be empty",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "",
			},
		}, {
			name: "Bad url (400)",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				wantError:   true,
				body:        "URL is not an URL format",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "!!!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, s)
			handle := http.HandlerFunc(handler.CreateShortURLText)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.want.wantError {
				require.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.isNew {
				assert.NotEmpty(t, string(response))
				assert.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Contains(t, string(response), tt.want.body)
		})
	}
}

func TestFail(t *testing.T) {
	c := testConfig()

	s := inmemory.NewInMemory()
	handler := NewHandler(c, s)

	type request struct {
		method string
		target string
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Serve any request and fail with 400",
			request: request{
				method: http.MethodGet,
				target: "/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "bad protocol",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Fail)
			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			require.Contains(t, string(response), tt.want.body)
		})
	}
}

func TestCreateShortURLJSON(t *testing.T) {
	c := testConfig()
	s := inmemory.NewInMemory()
	ctx := _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy")
	s.Store(_context.WithValue(ctx, context.ContextUserKey, "DoomGuy"), domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	})
	type want struct {
		contentType string
		statusCode  int
		isNew       bool
		wantError   bool
		body        string
	}
	type request struct {
		method string
		target string
		body   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Get old short url from full (409)",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusConflict,
				isNew:       false,
				wantError:   false,
				body:        strings.Join([]string{`{"result":"`, string(c.BaseURL), `/short"}`}, ""),
			},
			request: request{
				method: "POST",
				target: "/api/shorten",
				body:   `{"url":"http://www.yandex.ru/verylongpath"}`,
			},
		}, {
			name: "Create new short url from full (201)",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				isNew:       true,
				wantError:   false,
				body:        string(c.BaseURL),
			},
			request: request{
				method: "POST",
				target: "/api/shorten",
				body:   `{"url":"http://www.yandex.ru/very"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, s)
			handle := http.HandlerFunc(handler.CreateShortURLJSON)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.want.wantError {
				require.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.isNew {
				assert.NotEmpty(t, string(response))
				assert.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Contains(t, string(response), tt.want.body)
		})
	}
}

func TestHandler_CreateShortURLBatch(t *testing.T) {
	c := testConfig()
	s := inmemory.NewInMemory()
	ctx := _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy")
	s.StoreBatch(ctx, map[string]domain.URL{
		"1": {
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath",
			Short:  "short",
		},
	})
	type want struct {
		contentType string
		statusCode  int
		isNew       bool
		wantError   bool
		body        string
	}
	type request struct {
		method string
		target string
		body   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Batch create short url from full (201)",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				isNew:       false,
				wantError:   false,
				body:        `[{"correlation_id":"1","short_url":"` + string(c.BaseURL) + `/short"}]`,
			},
			request: request{
				method: "POST",
				target: "/api/shorten/batch",
				body:   `[{"correlation_id":"1","original_url":"http://www.yandex.ru/verylongpath"}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, s)
			handle := http.HandlerFunc(handler.CreateShortURLBatch)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.want.wantError {
				require.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.isNew {
				assert.NotEmpty(t, string(response))
				assert.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Contains(t, string(response), tt.want.body)
		})
	}
}

func TestHandler_GetUserURLs(t *testing.T) {
	c := testConfig()
	s := inmemory.NewInMemory()
	ctx := _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy")
	s.Store(ctx, domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	})
	type want struct {
		contentType string
		statusCode  int
		wantError   bool
		body        string
	}
	type request struct {
		method string
		target string
		ctx    _context.Context
		body   string
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Get current user URLs (200)",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				wantError:   false,
				body:        `[{"original_url":"http://www.yandex.ru/verylongpath","short_url":"` + string(c.BaseURL) + `/short"}]`,
			},
			request: request{
				method: "POST",
				target: "/api/user/urls",
				ctx:    _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy"),
				body:   ``,
			},
		},
		{
			name: "Get empty list user URLs (204)",
			want: want{
				contentType: "",
				statusCode:  http.StatusNoContent,
				wantError:   false,
				body:        ``,
			},
			request: request{
				method: "POST",
				target: "/api/user/urls",
				ctx:    _context.WithValue(_context.Background(), context.ContextUserKey, "Heretic"),
				body:   ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(tt.request.ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, s)
			handle := http.HandlerFunc(handler.GetUserURLs)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.want.wantError {
				require.Contains(t, string(response), tt.want.body)
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if len(tt.want.body) != 0 {
				assert.JSONEq(t, string(response), tt.want.body)
			}
		})
	}
}
