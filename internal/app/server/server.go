package server

import (
	"fmt"
	"net/http"

	"github.com/mikesvis/short/internal/app/helpers"
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

	helpers.SetURLOptions(myServerOptions.scheme, myServerOptions.host, myServerOptions.port)

	s = storage.NewStorageURL(make(map[string]domain.URL))
}

// Запуск сервера
func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`GET /`, func(w http.ResponseWriter, r *http.Request) {
		ServeGet(w, r, s)
	})
	mux.HandleFunc(`POST /`, func(w http.ResponseWriter, r *http.Request) {
		ServePost(w, r, s)
	})
	mux.HandleFunc(`/`, ServeOther)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", myServerOptions.host, myServerOptions.port), mux)
}
