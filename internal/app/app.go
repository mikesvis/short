package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	"github.com/mikesvis/short/internal/server"
	"github.com/mikesvis/short/internal/storage"
	"github.com/mikesvis/short/pkg/compressor"
	"go.uber.org/zap"
)

type App struct {
	config  *config.Config
	logger  *zap.SugaredLogger
	storage storage.Storage
	router  *chi.Mux
}

func New() *App {
	config := config.NewConfig()
	logger := logger.NewLogger()
	storage := storage.NewStorage(config)
	handler := server.NewHandler(config, storage)
	router := server.NewRouter(
		handler,
		middleware.RequestResponseLogger(logger),
		compressor.GZip(
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

func (a *App) Run() {
	a.logger.Infow("Config initialized", "config", a.config)
	// хоспади какие костыли ради разных хранилищ
	if _, isCloser := a.storage.(storage.StorageCloser); isCloser {
		defer a.storage.(storage.StorageCloser).Close()
	}
	if err := http.ListenAndServe(string(a.config.ServerAddress), a.router); err != nil {
		a.logger.Fatalw(err.Error(), "event", "start server")
	}
}
