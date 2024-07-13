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
	ctx := r.Context()
	shortKey := strings.TrimLeft(r.RequestURI, "/")
	item, err := h.storage.GetByShort(ctx, shortKey)
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
	ctx := r.Context()
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

	item, err := h.storage.GetByFull(ctx, URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if (item == domain.URL{}) {
		item = domain.URL{
			Full:  URL,
			Short: keygen.GetRandkey(keygen.KeyLength),
		}
		h.storage.Store(ctx, item)
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
	ctx := r.Context()
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

	item, err := h.storage.GetByFull(ctx, URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if (item == domain.URL{}) {
		item = domain.URL{
			Full:  URL,
			Short: keygen.GetRandkey(keygen.KeyLength),
		}
		h.storage.Store(ctx, item)
		status = http.StatusCreated
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := api.Response{Result: api.URL(urlformat.FormatURL(string(h.config.BaseURL), item.Short))}
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}

// Пинг хранилки
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.storage.Ping(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработка /api/shorten/batch POST
// Проверка на битый JSON
// Генерация коротких урлов пачкой
// Проверка на пустой URL
// Проверка на валидность URL
// Запись сокращенного URL в условную "базу" если нет такого ключа
func (h *Handler) CreateShortURLBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request []api.BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// генерим потенциальные domain.URL на сохранение с свежими Short
	pack := make(map[string]domain.URL)
	for _, v := range request {
		pack[string(v.CorrelationID)] = domain.URL{
			Full:  string(v.OriginalURL),
			Short: keygen.GetRandkey(keygen.KeyLength),
		}
	}

	// domain.URL.Short в процессе сохранения поменяем на старый если такой domain.URL.Full уже есть
	// цель: сделать получение/вычисление/сохранение в одну транзакцию
	// упарываться ответить кодом 201 если всех элементов не существовало ранее - не стал
	// TODO - не учитывается кейс с двумя одинаковыми correlation_id (не описаны требования!) :)
	// TODO поиск в хранилке надо сделать через кэш, мотать файл/базу каждый раз не комильфо
	stored, err := h.storage.StoreBatch(ctx, pack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response []api.BatchResponse
	for k, v := range stored {
		response = append(response, api.BatchResponse{
			CorrelationID: api.CorrelationID(k),
			ShortURL:      api.ShortURL(urlformat.FormatURL(string(h.config.BaseURL), v.Short)),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}
