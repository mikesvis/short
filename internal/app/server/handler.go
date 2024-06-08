package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mikesvis/short/internal/app/config"
	"github.com/mikesvis/short/internal/app/helpers"
	"github.com/mikesvis/short/internal/app/storage"
	"github.com/mikesvis/short/internal/domain"
)

// Обработка Get
// Получение короткого URL из запроса
// Поиск в условной "базе" полного URL по сокращенному
func ServeGet(s storage.StorageURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortKey := strings.TrimLeft(r.RequestURI, "/")
		item := s.GetByShort(shortKey)
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

func getScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}

// Обработка POST
// Проверка на пустое тело запроса
// Проверка на валидность URL
// Запись сокращенного Url в условную "базу" если нет такого ключа
func ServePost(s storage.StorageURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = helpers.ValidateURL(string(body))
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		URL := helpers.SanitizeURL(string(body))
		status := http.StatusOK

		item := s.GetByFull(URL)
		if s.GetByFull(URL) == (domain.URL{}) {
			item = domain.URL{
				Full:  URL,
				Short: helpers.GetRandkey(helpers.KeyLength),
			}
			s.Store(item)
			status = http.StatusCreated

			log.Printf("add new item to urls %+v", item)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(status)
		w.Write([]byte(helpers.FormatURL(config.GetShortLinkAddr(), item.Short)))
	}
}

// Обработка всего остального
func ServeOther(w http.ResponseWriter, r *http.Request) {
	err := errors.New("bad protocol")
	log.Printf("%s", err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
