package server

import (
	_context "context"
	goerrors "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/storage"
	mock_storage "github.com/mikesvis/short/mocks/storage"
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
	ctx, cancel := _context.WithCancel(_context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)

	gomock.InOrder(
		mockedStorage.EXPECT().GetByShort(ctx, "short1").Return(domain.URL{
			UserID:  "DoomGuy",
			Full:    "http://www.yandex.ru/verylongpath",
			Short:   "short",
			Deleted: false,
		}, nil),
		mockedStorage.EXPECT().GetByShort(ctx, "short2").Return(domain.URL{}, goerrors.New("full url is not found")),
		mockedStorage.EXPECT().GetByShort(ctx, "short3").Return(domain.URL{
			UserID:  "DoomGuy",
			Full:    "http://www.yandex.ru/verylongpath",
			Short:   "short3",
			Deleted: true,
		}, nil),
	)

	type want struct {
		statusCode  int
		newLocation string
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
			},
			request: request{
				methhod: "GET",
				target:  "/short1",
			},
		}, {
			name: "Full url does not exist (400)",
			want: want{
				statusCode:  http.StatusBadRequest,
				newLocation: "",
				body:        "full url is not found",
			},
			request: request{
				methhod: "GET",
				target:  "/short2",
			},
		}, {
			name: "Item has gone",
			want: want{
				statusCode:  http.StatusGone,
				newLocation: "",
			},
			request: request{
				methhod: "GET",
				target:  "/short3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.methhod, tt.request.target, nil)
			w := httptest.NewRecorder()
			handler := NewHandler(c, mockedStorage)
			handle := http.HandlerFunc(handler.GetFullURL)
			handle(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if len(tt.want.newLocation) > 0 {
				assert.Contains(t, result.Header.Values("Location"), tt.want.newLocation)
			}

			response, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			if len(tt.want.body) > 0 {
				require.Contains(t, string(response), tt.want.body)
			}
		})
	}
}

func BenchmarkGetFullURL(b *testing.B) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	s.Store(_context.Background(), domain.URL{
		Full:  "http://www.yandex.ru/verylongpath",
		Short: "short",
	})

	request := httptest.NewRequest("GET", "/short", nil)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.GetFullURL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(w, request)
	}
}

func TestCreateShortURLText(t *testing.T) {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	ctxMock, cancel := _context.WithCancel(_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"))
	defer cancel()
	ctxReq := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)

	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("jHQri").Times(1)
	mockedStorage.EXPECT().Store(ctxMock, gomock.Any()).Return(domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "jHQri",
	}, nil).Times(1)
	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("bHgdb").Times(1)
	mockedStorage.EXPECT().Store(ctxMock, gomock.Any()).Return(domain.URL{}, goerrors.New("Fail")).Times(1)

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
			name: "Create new short url from full (201)",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				isNew:       true,
				body:        string(c.BaseURL),
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/verylongpath",
			},
		},
		{
			name: "Something is rotten (400)",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				body:        "Fail",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "http://www.yandex.ru/verylongpath2",
			},
		},
		{
			name: "Empty body (400)",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
				body:        "URL can not be empty",
			},
			request: request{
				method: "POST",
				target: "/",
				body:   "",
			},
		},
		{
			name: "Bad url (400)",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				isNew:       false,
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
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctxReq)
			w := httptest.NewRecorder()
			handler := NewHandler(c, mockedStorage)
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

func BenchmarkCreateShortURLText(b *testing.B) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	s.Store(ctx, domain.URL{
		Full:   "http://www.yandex.ru/verylongpath",
		UserID: "Doomguy",
		Short:  "short",
	})

	request := httptest.NewRequest("POST", "/", strings.NewReader("http://www.yandex.ru/verylongpath")).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(w, request)
	}
}

func TestFail(t *testing.T) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
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
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	s.Store(_context.WithValue(ctx, context.UserIDContextKey, "DoomGuy"), domain.URL{
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

func BenchmarkCreateShortURLJSON(b *testing.B) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	s.Store(_context.WithValue(ctx, context.UserIDContextKey, "DoomGuy"), domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	})

	request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://www.yandex.ru/verylongpath"}`)).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLText)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(w, request)
	}
}

func TestCreateShortURLBatch(t *testing.T) {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	ctxMock, cancel := _context.WithCancel(_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"))
	defer cancel()
	ctxReq := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)

	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("short1")
	mockedStorage.EXPECT().StoreBatch(ctxMock, map[string]domain.URL{
		"1": {
			UserID:  "DoomGuy",
			Full:    "http://www.yandex.ru/verylongpath",
			Short:   "short1",
			Deleted: false,
		},
	}).Return(map[string]domain.URL{
		"1": {
			UserID:  "DoomGuy",
			Full:    "http://www.yandex.ru/verylongpath",
			Short:   "short1",
			Deleted: false,
		},
	}, nil)

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
				body:        `[{"correlation_id":"1","short_url":"` + string(c.BaseURL) + `/short1"}]`,
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
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctxReq)
			w := httptest.NewRecorder()
			handler := NewHandler(c, mockedStorage)
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

func BenchmarkCreateShortURLBatch(b *testing.B) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	s.StoreBatch(ctx, map[string]domain.URL{
		"1": {
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath",
			Short:  "short",
		},
	})

	request := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(`[{"correlation_id":"1","original_url":"http://www.yandex.ru/verylongpath"}]`)).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLBatch)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(w, request)
	}
}

func TestGetUserURLs(t *testing.T) {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)

	mockedStorage.EXPECT().GetUserURLs(gomock.Any(), gomock.Eq("DoomGuy")).Return([]domain.URL{
		{
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath",
			Short:  "short",
		},
	}, nil)

	mockedStorage.EXPECT().GetUserURLs(gomock.Any(), gomock.Eq("Heretic")).Return([]domain.URL{}, nil)

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
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
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
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "Heretic"),
				body:   ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(tt.request.ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, mockedStorage)
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

func BenchmarkGetUserURLs(b *testing.B) {
	c := testConfig()
	l, _ := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	s.Store(ctx, domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	})

	request := httptest.NewRequest("POST", "/api/user/urls", strings.NewReader(``)).WithContext(_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"))
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.GetUserURLs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(w, request)
	}
}

func TestHandler_Ping(t *testing.T) {
	l, _ := logger.NewLogger()
	ctx := _context.Background()
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()
	type args struct {
		config *config.Config
	}
	type want struct {
		statusCode int
		wantError  bool
	}
	tests := []struct {
		args args
		name string
		want want
	}{
		{
			name: "In memory db not found",
			args: args{
				config: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "",
				},
			},
			want: want{
				wantError:  false,
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Ping file db success",
			args: args{
				config: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: tmpFile.Name(),
					DatabaseDSN:     "",
				},
			},
			want: want{
				wantError:  false,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Ping file db fail",
			args: args{
				config: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "dummyfile.bin",
					DatabaseDSN:     "",
				},
			},
			want: want{
				wantError:  false,
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "Postgres db fail",
			args: args{
				config: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "host=0.0.0.0 port=5432 user=postgres password=postgres dbname=nodb sslmode=disable",
				},
			},
			want: want{
				wantError:  true,
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/ping", strings.NewReader(``)).WithContext(ctx)
			w := httptest.NewRecorder()
			s, err := storage.NewStorage(tt.args.config, l)
			if tt.want.wantError {
				require.Error(t, err)
			}
			handler := NewHandler(tt.args.config, s)
			handle := http.HandlerFunc(handler.Ping)
			handle(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
	os.Remove(tmpFile.Name())
}

func TestHander_DeleteUserURLs(t *testing.T) {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	ctxMock, cancel := _context.WithCancel(_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"))
	defer cancel()
	ctxReq := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)
	mockedStorage.EXPECT().DeleteBatch(ctxMock, "DoomGuy", []string{"short1", "short2"}).Return()

	type want struct {
		statusCode int
	}
	type request struct {
		method string
		target string
		body   string
	}
	tests := []struct {
		name    string
		want    want
		arg     storage.Storage
		request request
	}{
		{
			name: "Delete user URLs Success",
			arg:  mockedStorage,
			want: want{
				statusCode: http.StatusAccepted,
			},
			request: request{
				method: "DELETE",
				target: "/api/user/urls",
				body:   `["short1","short2"]`,
			},
		},
		{
			name: "Delete user URLs Bad Request",
			arg:  mockedStorage,
			want: want{
				statusCode: http.StatusBadRequest,
			},
			request: request{
				method: "DELETE",
				target: "/api/user/urls",
				body:   ``,
			},
		},
		{
			name: "Delete user URLs Bad Request empty request",
			arg:  mockedStorage,
			want: want{
				statusCode: http.StatusBadRequest,
			},
			request: request{
				method: "DELETE",
				target: "/api/user/urls",
				body:   `[]`,
			},
		},
		{
			name: "Storage doesnt support batch delete",
			arg:  &inmemory.InMemory{},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
			request: request{
				method: "DELETE",
				target: "/api/user/urls",
				body:   `["short1","short2"]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(ctxReq)
			w := httptest.NewRecorder()
			handler := NewHandler(c, tt.arg)
			handle := http.HandlerFunc(handler.DeleteUserURLs)
			handle(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestGetStats(t *testing.T) {
	c := &config.Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStorage := mock_storage.NewMockStorageDeleter(ctrl)

	mockedStorage.EXPECT().GetStats(gomock.Any()).Return(domain.Stats{
		URLs:  0,
		Users: 0,
	}, nil)

	mockedStorage.EXPECT().GetStats(gomock.Any()).Return(domain.Stats{}, goerrors.New("Failed"))

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
			name: "Get zero stats",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				wantError:   false,
				body:        `{"urls":0,"users":0}`,
			},
			request: request{
				method: "GET",
				target: "/api/internal/urls",
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				body:   ``,
			},
		},
		{
			name: "Get error",
			want: want{
				contentType: "",
				statusCode:  http.StatusInternalServerError,
				wantError:   true,
				body:        ``,
			},
			request: request{
				method: "POST",
				target: "/api/user/urls",
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "Heretic"),
				body:   ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, strings.NewReader(tt.request.body)).WithContext(tt.request.ctx)
			w := httptest.NewRecorder()
			handler := NewHandler(c, mockedStorage)
			handle := http.HandlerFunc(handler.GetStats)
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
