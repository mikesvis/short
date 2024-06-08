package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/app/storage"
	"github.com/mikesvis/short/internal/domain"
)

type serverOptions struct {
	scheme,
	host,
	port string
}

var myServerOptions serverOptions

var s storage.StorageURL

func init() {
	myServerOptions = serverOptions{
		"http",
		"localhost",
		"8080",
	}

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
	return http.ListenAndServe(fmt.Sprintf("%s:%s", myServerOptions.host, myServerOptions.port), ShortRouter())
}
