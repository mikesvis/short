package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/middleware"
)

func NewRouter(h *Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middlewares...)

	r.Route("/api", func(r chi.Router) {
		r.With(middleware.SignIn).Post("/shorten/batch", h.CreateShortURLBatch)
		r.With(middleware.SignIn).Post("/shorten", h.CreateShortURLJSON)
		r.With(middleware.Auth).Get("/user/urls", h.GetUserURLs)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", h.Ping)
		r.Get("/{shortKey}", h.GetFullURL)
		r.With(middleware.SignIn).Post("/", h.CreateShortURLText)
		r.Get("/", h.Fail)
		r.Patch("/", h.Fail)
		r.Put("/", h.Fail)
		r.Delete("/", h.Fail)
	})

	return r
}
