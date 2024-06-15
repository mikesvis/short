package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/storage"
)

var s StorageURL

func init() {
	s = storage.NewStorageURL(make(map[domain.ID]domain.URL))
}

func NewRouter() *chi.Mux {
	r := chi.NewMux()
	h := NewHandler(s)
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortKey}", h.ServeGet())
		r.Post("/", h.ServePost())
		r.Get("/", h.ServeOther)
		r.Patch("/", h.ServeOther)
		r.Put("/", h.ServeOther)
		r.Delete("/", h.ServeOther)
	})

	return r
}

// Запуск сервера
func Run() error {
	return http.ListenAndServe(config.GetServerAddress(), NewRouter())
}
