package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mikesvis/short/internal/app/helpers"
	"github.com/mikesvis/short/internal/app/storage"
	"github.com/mikesvis/short/internal/domain"
)

// Обработка GET
func ServeGet(w http.ResponseWriter, r *http.Request, s storage.StorageURL) {
	shortURL := fmt.Sprintf("%s://%s%s", getScheme(r), r.Host, r.URL.Path)
	item := s.GetByShort(shortURL)
	if item == (domain.URL{}) {
		errorResponse(w, fmt.Errorf("full url is not found for %s", shortURL))
		return
	}

	w.Header().Set("Location", item.Full)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}

// Обработка POST
func ServePost(w http.ResponseWriter, r *http.Request, s storage.StorageURL) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = helpers.ValidateURL(string(body))
	if err != nil {
		errorResponse(w, err)
		return
	}

	URL := helpers.SanitizeURL(string(body))
	item := s.GetByFull(URL)
	if s.GetByFull(URL) == (domain.URL{}) {
		item = domain.URL{
			Full:  URL,
			Short: helpers.GetFormattedURL(helpers.GetRandkey(helpers.KeyLength)),
		}
		s.Store(item)
		log.Printf("add new item to urls %+v", item)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(item.Short))
}

// Получение Url с сокращением если есть в условной "базе", генерация нового если в "базе" нет
// func getShortURL(urlToShort string) string {
// 	if shortKey := findKeyByValueInMap(urlToShort, storage); shortKey != "" {
// 		return getFormattedURL(shortKey)
// 	}

// 	shortKey := getRandkey(keyLength)
// 	storage[shortKey] = urlToShort
// 	shortURL := helpers.GetFormattedURL(shortKey)

// 	log.Printf("generated and saved new shorten key for %s: %s", urlToShort, shortKey)

// 	return shortURL
// }

// Обработка всего остального
func ServeOther(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, errors.New("bad protocol"))
}

func errorResponse(w http.ResponseWriter, err error) {
	log.Printf("%s", err)
	w.WriteHeader(http.StatusBadRequest)
}
