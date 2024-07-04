package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/compressor"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/storage"
)

func NewRouter(c *config.Config, s storage.Storage, h *Handler) *chi.Mux {
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
