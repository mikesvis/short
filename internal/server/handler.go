package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mikesvis/short/internal/api"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/keygen"
	"github.com/mikesvis/short/internal/storage"
	"github.com/mikesvis/short/pkg/urlformat"
)

type Handler struct {
	config  *config.Config
	storage storage.Storage
}

func NewHandler(config *config.Config, storage storage.Storage) *Handler {
	return &Handler{config, storage}
}

// Обработка Get
// Получение короткого URL из запроса
// Поиск в условной "базе" полного URL по сокращенному
func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimLeft(r.RequestURI, "/")
	item, err := h.storage.GetByShort(shortKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if (item == domain.URL{}) {
		err := fmt.Errorf("full url is not found for %s", shortKey)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Location", item.Full)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Обработка POST
// Проверка на пустое тело запроса
// Проверка на валидность URL
// Запись сокращенного Url в условную "базу" если нет такого ключа
func (h *Handler) CreateShortURLText(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = urlformat.ValidateURL(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	URL := urlformat.SanitizeURL(string(body))
	status := http.StatusOK

	item, err := h.storage.GetByFull(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if (item == domain.URL{}) {
		item = domain.URL{
			Full:  URL,
			Short: keygen.GetRandkey(keygen.KeyLength),
		}
		h.storage.Store(item)
		status = http.StatusCreated
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(urlformat.FormatURL(string(h.config.BaseURL), item.Short)))
}

// Обработка всего остального
func (h *Handler) Fail(w http.ResponseWriter, r *http.Request) {
	err := errors.New("bad protocol")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// Обработка /api/shorten POST
// Проверка на битый JSON
// Проверка на пустой URL
// Проверка на валидность URL
// Запись сокращенного URL в условную "базу" если нет такого ключа
func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	var request api.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	URL := string(request.URL)
	err := urlformat.ValidateURL(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	URL = urlformat.SanitizeURL(URL)
	status := http.StatusOK

	item, err := h.storage.GetByFull(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if (item == domain.URL{}) {
		item = domain.URL{
			Full:  URL,
			Short: keygen.GetRandkey(keygen.KeyLength),
		}
		h.storage.Store(item)
		status = http.StatusCreated
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := api.Response{Result: api.URL(urlformat.FormatURL(string(h.config.BaseURL), item.Short))}
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}
