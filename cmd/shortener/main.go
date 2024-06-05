package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const keyLength = 5

type serverOptions struct {
	scheme,
	host,
	port string
}

var myServerOptions serverOptions
var storage map[string]string

func init() {
	myServerOptions = serverOptions{
		"http",
		"127.0.0.1",
		"8080",
	}
	storage = make(map[string]string)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// Запуск сервера
func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, urlProcessor)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", myServerOptions.host, myServerOptions.port), mux)
}

// Обработчик запросов
func urlProcessor(w http.ResponseWriter, r *http.Request) {

	// вот тут я помучался поскольку в тестах при POST Content-Type приходит еще кодировка
	// а в GET не приходит заголовок Content-Type вообще
	// в задании же указано Content-Type: text/plain в POST и GET :) happy debugging!
	//
	// contentType := r.Header.Get("Content-Type")
	// if !strings.Contains(contentType, "text/plain") {
	// 	log.Printf("Content-Type %s is not allowed\r\n", contentType)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	switch r.Method {
	case http.MethodGet:
		serveGet(w, r)
		return
	case http.MethodPost:
		servePost(w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

// Обработка GET
func serveGet(w http.ResponseWriter, r *http.Request) {
	fullUrl, err := getFullUrl(r)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fullUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Получение полного Url если есть соответствие в условной "базе"
func getFullUrl(r *http.Request) (string, error) {
	shortKeyFromUrl := strings.TrimLeft(string(r.URL.Path), "/")
	if fullUrl := findValueByKeyInMap(shortKeyFromUrl, storage); fullUrl != "" {
		return fullUrl, nil

	}

	return "", fmt.Errorf("short url not found for key: %s", shortKeyFromUrl)
}

// Поиск значения по ключу
func findValueByKeyInMap(needle string, storage map[string]string) string {
	for k, v := range storage {
		if k != needle {
			continue
		}

		log.Printf("short %s key is found for url %s", k, v)
		return v
	}

	return ""
}

// Обработка POST
func servePost(w http.ResponseWriter, r *http.Request) {
	body, err := getUrlToShort(r.Body)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortUrl := getShortUrl(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

// Получение и проверка валидности Url
func getUrlToShort(buffer io.ReadCloser) (string, error) {
	body, err := io.ReadAll(buffer)
	if err != nil {
		return "", fmt.Errorf("error while reading POST body: %w", err)
	}

	if len(body) == 0 {
		return "", errors.New("POST body can not be empty")
	}

	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		return "", fmt.Errorf("POST body is not an URL format, %s given", err)
	}

	return strings.Trim(string(body), ""), nil
}

// Получение Url с сокращением если есть в условной "базе", генерация нового если в "базе" нет
func getShortUrl(urlToShort string) string {
	if shortKey := findKeyByValueInMap(urlToShort, storage); shortKey != "" {
		return getFormattedUrl(shortKey)
	}

	shortKey := getRandkey(keyLength)
	storage[shortKey] = urlToShort
	shortUrl := getFormattedUrl(shortKey)

	log.Printf("generated and saved new shorten key for %s: %s", urlToShort, shortKey)

	return shortUrl
}

// Поиск ключа по значению
func findKeyByValueInMap(needle string, storage map[string]string) string {
	for k, v := range storage {
		if v != needle {
			continue
		}

		log.Printf("short key was generated already for %s: %s", v, k)
		return k
	}

	return ""
}

// Шаблон сокращенного Url
func getFormattedUrl(shortKey string) string {
	return fmt.Sprintf("%s://%s:%s/%s", myServerOptions.scheme, myServerOptions.host, myServerOptions.port, shortKey)
}

// Получение рандомного ключа/строки
func getRandkey(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
