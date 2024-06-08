package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/app/config"
	"github.com/mikesvis/short/internal/app/storage"
	"github.com/mikesvis/short/internal/domain"
)

var s storage.StorageURL

func init() {
	s = storage.NewStorageURL(make(map[domain.ID]domain.URL))
}

func ShortRouter() chi.Router {
	r := chi.NewRouter()
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
	return http.ListenAndServe(config.GetServerHostAddr(), ShortRouter())
}
