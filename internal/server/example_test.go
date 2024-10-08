package server

import (
	_context "context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/logger"
)

func ExampleHandler_GetFullURL() {
	// определяем конфиг, логгер, хранилку (в примере хранилка в памяти)
	c := testConfig()
	l := logger.NewLogger()
	s := inmemory.NewInMemory(l)

	// формируем запрос
	request := httptest.NewRequest("GET", "http://example.com/short", nil)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.GetFullURL)

	// отправляем запрос и получаем результат
	handle(w, request)
	result := w.Result()
	defer result.Body.Close()
	response, _ := io.ReadAll(result.Body)
	fmt.Printf("%v", response)
}

func ExampleHandler_CreateShortURLText() {
	// определяем конфиг, логгер, хранилку (в примере хранилка в памяти) и конекст с ID пользователя
	c := testConfig()
	l := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	// формируем запрос
	request := httptest.NewRequest("POST", "/", strings.NewReader("http://www.yandex.ru/verylongpath")).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLText)

	// отправляем запрос и получаем результат
	handle(w, request)
	result := w.Result()
	defer result.Body.Close()
	response, _ := io.ReadAll(result.Body)

	fmt.Printf("%v", response)
}

func ExampleHandler_CreateShortURLJSON() {
	// определяем конфиг, логгер, хранилку (в примере хранилка в памяти) и конекст с ID пользователя
	c := testConfig()
	l := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	// формируем запрос
	request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://www.yandex.ru/verylongpath"}`)).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLText)

	// отправляем запрос и получаем результат
	handle(w, request)
	result := w.Result()
	defer result.Body.Close()
	response, _ := io.ReadAll(result.Body)

	fmt.Printf("%v", response)
}

func ExampleHandler_CreateShortURLBatch() {
	// определяем конфиг, логгер, хранилку (в примере хранилка в памяти) и конекст с ID пользователя
	c := testConfig()
	l := logger.NewLogger()
	s := inmemory.NewInMemory(l)
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	// формируем запрос
	request := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(`[{"correlation_id":"1","original_url":"http://www.yandex.ru/verylongpath"}]`)).WithContext(ctx)
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.CreateShortURLBatch)

	// отправляем запрос и получаем результат
	handle(w, request)
	result := w.Result()
	defer result.Body.Close()
	response, _ := io.ReadAll(result.Body)

	fmt.Printf("%v", response)
}

func ExampleHandler_GetUserURLs() {
	// определяем конфиг, логгер, хранилку (в примере хранилка в памяти) и конекст с ID пользователя
	c := testConfig()
	l := logger.NewLogger()
	s := inmemory.NewInMemory(l)

	// формируем запрос
	request := httptest.NewRequest("POST", "/api/user/urls", strings.NewReader(``)).WithContext(_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"))
	w := httptest.NewRecorder()
	handler := NewHandler(c, s)
	handle := http.HandlerFunc(handler.GetUserURLs)

	// отправляем запрос и получаем результат
	handle(w, request)
	result := w.Result()
	defer result.Body.Close()
	response, _ := io.ReadAll(result.Body)

	fmt.Printf("%v", response)
}
