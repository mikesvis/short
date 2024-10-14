// Пакет приложения сокращателя ссылок.
package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	"github.com/mikesvis/short/internal/server"
	"github.com/mikesvis/short/internal/storage"
	"go.uber.org/zap"
)

// App - стуктура приложения с конфигом, логгером, storage и роутером.
type App struct {
	config  *config.Config
	logger  *zap.SugaredLogger
	storage storage.Storage
	router  *chi.Mux
}

// Конструктор приложения, здесь инициализируются все зависимости:
// конфиг приложения, логгер, storage, роутер. Также здесь регистрируются middleware приложения.
func New(config *config.Config) *App {
	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	storage, err := storage.NewStorage(config, logger)
	if err != nil {
		panic(err)
	}

	handler := server.NewHandler(config, storage)
	router := server.NewRouter(
		handler,
		middleware.RequestResponseLogger(logger),
		middleware.GZip(
			[]string{
				"application/json",
				"text/html",
				"application/x-gzip",
			},
		),
	)

	return &App{
		config,
		logger,
		storage,
		router,
	}
}

// Запуск приложения.
func (a *App) Run() error {
	a.logger.Infow("Config initialized", "config", a.config)
	if _, isCloser := a.storage.(storage.StorageCloser); isCloser {
		defer a.storage.(storage.StorageCloser).Close()
	}

	var err error
	if a.config.EnableHTTPS {
		err = http.ListenAndServeTLS(a.config.ServerAddress, a.config.ServerCertPath, a.config.ServerKeyPath, a.router)
	} else {
		err = http.ListenAndServe(a.config.ServerAddress, a.router)
	}

	if err != nil {
		a.logger.Errorf(err.Error(), "event", "start server")
		return err
	}

	return nil
}
