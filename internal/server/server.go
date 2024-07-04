package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/compressor"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/storage/filedb"
	"github.com/mikesvis/short/internal/storage/memorymap"
)

var s StorageURL

func NewRouter() *chi.Mux {
	s = newStorage(config.GetFileStoragePath())
	h := NewHandler(s)

	r := chi.NewMux()
	r.Use(logger.RequestResponseLogger)
	r.Use(compressor.GZip)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.CreateShortURLJSON)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortKey}", h.GetFullURL)
		r.Post("/", h.CreateShortURLText)
		r.Get("/", h.Fail)
		r.Patch("/", h.Fail)
		r.Put("/", h.Fail)
		r.Delete("/", h.Fail)
	})

	return r
}

func newStorage(fileStoragePath string) StorageURL {
	if len(fileStoragePath) == 0 {
		logger.Log.Info("Using in-memory map storage")
		return memorymap.NewStorageURL(make(map[domain.ID]domain.URL))
	}

	logger.Log.Infof("Using file storage by path %s", fileStoragePath)
	return filedb.NewStorageURL(fileStoragePath)
}

// Запуск сервера
func Run() error {
	logger.Log.Infow("Running server", "address", config.GetServerAddress())
	return http.ListenAndServe(config.GetServerAddress(), NewRouter())
}
