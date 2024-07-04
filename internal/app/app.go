package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/server"
	"github.com/mikesvis/short/internal/storage"
)

type App struct {
	config  *config.Config
	storage storage.Storage
	router  *chi.Mux
}

func New() *App {
	config := config.NewConfig()
	storage := storage.NewStorage(string(config.FileStoragePath))
	handler := server.NewHandler(config, storage)
	router := server.NewRouter(config, storage, handler)
	return &App{
		config,
		storage,
		router,
	}
}

func (a *App) Run() error {
	logger.Log.Infow("Config initialized", "config", a.config)
	return http.ListenAndServe(string(a.config.ServerAddress), a.router)
}
