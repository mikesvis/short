package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middlewares...)

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
