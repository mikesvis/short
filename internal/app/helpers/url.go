package helpers

import (
	"fmt"
	"net/url"
	"strings"
)

type URLOptions struct {
	scheme,
	host,
	port string
}

var myURLOptions URLOptions

func SetURLOptions(sheme, host, port string) {
	myURLOptions = URLOptions{
		scheme: sheme,
		host:   host,
		port:   port,
	}
}

// Шаблон сокращенного Url
func GetFormattedURL(shortKey string) string {
	return fmt.Sprintf("%s://%s:%s/%s", myURLOptions.scheme, myURLOptions.host, myURLOptions.port, shortKey)
}

// Получение и проверка валидности Url
func ValidateURL(URL string) error {
	if len(URL) == 0 {
		return fmt.Errorf("POST body can not be empty")
	}

	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return fmt.Errorf("POST body is not an URL format, %s given", err)
	}

	return nil
}

func SanitizeURL(URL string) string {
	return strings.Trim(URL, "")
}
