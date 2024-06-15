package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/keygen"
	"github.com/mikesvis/short/internal/storage"
	"github.com/mikesvis/short/internal/urlformat"
)

type handler struct {
	storage storage.StorageURL
}

func NewHandler(s storage.StorageURL) *handler {
	return &handler{storage: s}
}

// Обработка Get
// Получение короткого URL из запроса
// Поиск в условной "базе" полного URL по сокращенному
func (h *handler) ServeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortKey := strings.TrimLeft(r.RequestURI, "/")
		item := h.storage.GetByShort(shortKey)
		if item == (domain.URL{}) {
			err := fmt.Errorf("full url is not found for %s", shortKey)
			http.Error(w, err.Error(), http.StatusBadRequest)

			log.Printf("%s", err)
			return
		}

		w.Header().Set("Location", item.Full)
		w.WriteHeader(http.StatusTemporaryRedirect)

		log.Printf("found item for full url %+v", item)
	}
}

// Обработка POST
// Проверка на пустое тело запроса
// Проверка на валидность URL
// Запись сокращенного Url в условную "базу" если нет такого ключа
func (h *handler) ServePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = urlformat.ValidateURL(string(body))
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		URL := urlformat.SanitizeURL(string(body))
		status := http.StatusOK

		item := h.storage.GetByFull(URL)
		if h.storage.GetByFull(URL) == (domain.URL{}) {
			item = domain.URL{
				Full:  URL,
				Short: keygen.GetRandkey(keygen.KeyLength),
			}
			h.storage.Store(item)
			status = http.StatusCreated

			log.Printf("add new item to urls %+v", item)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(status)
		w.Write([]byte(urlformat.FormatURL(config.GetBaseURL(), item.Short)))
	}
}

// Обработка всего остального
func (h *handler) ServeOther(w http.ResponseWriter, r *http.Request) {
	err := errors.New("bad protocol")
	log.Printf("%s", err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
