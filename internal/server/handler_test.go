package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeGet(t *testing.T) {
	type args struct {
		s StorageURL
	}
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
		args    args
		want    want
		request request
	}{
		{
			name: "Find full url (307)",
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath",
						Short: "short",
					},
				}),
			},
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
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{}),
			},
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
			handler := NewHandler(tt.args.s)
			handle := http.HandlerFunc(handler.ServeGet())
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

func TestServePost(t *testing.T) {
	config.InitConfig()
	type args struct {
		s StorageURL
	}
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
		args    args
		want    want
		request request
	}{
		{
			name: "Get short url from full (200)",
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath",
						Short: "short",
					},
				}),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
				isNew:       false,
				wantError:   false,
				body:        config.GetBaseURL() + "/short",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/verylongpath",
			},
		}, {
			name: "Create short url from full (201)",
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{}),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				isNew:       true,
				wantError:   false,
				body:        config.GetBaseURL(),
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/verylongpath",
			},
		}, {
			name: "Empty body (400)",
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{}),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				wantError:   true,
				body:        "",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "POST body can not be empty",
			},
		}, {
			name: "Bad url (400)",
			args: args{
				s: storage.NewStorageURL(map[domain.ID]domain.URL{}),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				wantError:   true,
				body:        "",
			},
			request: request{
				method: "POST",
				target: "/ya.ru",
				body:   "POST body is not an URL format",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			handler := NewHandler(tt.args.s)
			handle := http.HandlerFunc(handler.ServePost())
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

func TestServeOther(t *testing.T) {
	s := storage.NewStorageURL(map[domain.ID]domain.URL{})
	handler := NewHandler(s)

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
			h := http.HandlerFunc(handler.ServeOther)
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
