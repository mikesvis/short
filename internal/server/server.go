package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/storage"
)

var s storage.StorageURL

func init() {
	s = storage.NewStorageURL(make(map[domain.ID]domain.URL))
}

func NewRouter() *chi.Mux {
	r := chi.NewMux()
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortKey}", ServeGet(s))
		r.Post("/", ServePost(s))
		r.Get("/", ServeOther)
		r.Patch("/", ServeOther)
		r.Put("/", ServeOther)
		r.Delete("/", ServeOther)
	})

	return r
}

// Запуск сервера
func Run() error {
	return http.ListenAndServe(config.GetServerAddress(), NewRouter())
}
