// Модуль описания handler'ов.
package server

import (
	_context "context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/mikesvis/short/internal/api"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/keygen"
	"github.com/mikesvis/short/internal/storage"
	"github.com/mikesvis/short/pkg/urlformat"
)

// Хендлер приложения, включает в себя *config.Config и storage.Storage.
type Handler struct {
	config  *config.Config
	storage storage.Storage
}

// Конструктор хендлера
func NewHandler(config *config.Config, storage storage.Storage) *Handler {
	return &Handler{config, storage}
}

// Обработка Get
// Получение короткого URL из запроса
// Поиск в условной "базе" полного URL по сокращенному
func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

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

	if item.Deleted {
		w.WriteHeader(http.StatusGone)

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
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

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
	item := domain.URL{
		UserID: ctx.Value(context.UserIDContextKey).(string),
		Full:   URL,
		Short:  h.storage.GetRandkey(keygen.KeyLength),
	}
	status := http.StatusConflict

	item, err = h.storage.Store(ctx, item)
	if err != nil && !_errors.Is(err, errors.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err == nil {
		status = http.StatusCreated
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(urlformat.FormatURL(string(h.config.BaseURL), item.Short)))
}

// Обработка всего остального
func (h *Handler) Fail(w http.ResponseWriter, r *http.Request) {
	err := _errors.New("bad protocol")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// Обработка /api/shorten POST
// Проверка на битый JSON
// Проверка на пустой URL
// Проверка на валидность URL
// Запись сокращенного URL в условную "базу" если нет такого ключа
func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

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
	item := domain.URL{
		UserID: ctx.Value(context.UserIDContextKey).(string),
		Full:   URL,
		Short:  h.storage.GetRandkey(keygen.KeyLength),
	}
	status := http.StatusConflict

	item, err = h.storage.Store(ctx, item)
	if err != nil && !_errors.Is(err, errors.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err == nil {
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
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	if _, ok := h.storage.(storage.StoragePinger); !ok {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	err := h.storage.(storage.StoragePinger).Ping(ctx)

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
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	var request api.BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// генерим потенциальные domain.URL на сохранение с новым Short
	pack := make(map[string]domain.URL)
	for _, v := range request {
		pack[string(v.CorrelationID)] = domain.URL{
			UserID: ctx.Value(context.UserIDContextKey).(string),
			Full:   string(v.OriginalURL),
			Short:  h.storage.GetRandkey(keygen.KeyLength),
		}
	}

	// domain.URL.Short в процессе сохранения поменяем на старый если такой domain.URL.Full уже есть
	stored, err := h.storage.StoreBatch(ctx, pack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make(api.BatchResponse, 0, len(stored))
	for k, v := range stored {
		// не понимаю что мы тут сократили, по моему с BatchResponseItem было лучше (но исправил по замечанию ревью)
		response = append(response, struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: k,
			ShortURL:      urlformat.FormatURL(string(h.config.BaseURL), v.Short),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}

// Обработка /api/user/urls GET
// Получение URL пользователя
func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	// тут умышленно ctx, ctx.Value
	// 1ый аргумент - контекст, 2ой аргумент - само значение ID пользователя
	items, err := h.storage.GetUserURLs(ctx, ctx.Value(context.UserIDContextKey).(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make(api.UserResponse, 0, len(items))
	for _, v := range items {
		// не понимаю что мы тут сократили, по моему с UserResponseItem было лучше (но исправил по замечанию ревью)
		response = append(response, struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShortURL:    urlformat.FormatURL(string(h.config.BaseURL), v.Short),
			OriginalURL: v.Full,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(response)
}

// Обработка /api/user/urls DELETE
// Удаление URL пользователя
func (h *Handler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := _context.WithCancel(r.Context())
	defer cancel()

	if _, isDeleter := h.storage.(storage.StorageDeleter); !isDeleter {
		http.Error(w, fmt.Sprintf(`Batch delete is not supported for storage of type %s`, reflect.TypeOf(h.storage).String()), http.StatusInternalServerError)

		return
	}

	var request api.BatchDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// пачки нет
	if len(request) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.storage.(storage.StorageDeleter).DeleteBatch(ctx, ctx.Value(context.UserIDContextKey).(string), []string(request))

	w.WriteHeader(http.StatusAccepted)
}
