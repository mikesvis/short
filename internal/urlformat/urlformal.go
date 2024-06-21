package urlformat

import (
	"fmt"
	"net/url"
	"strings"
)

// Шаблон сокращенного Url
func FormatURL(linkServerAddress, shortKey string) string {
	return fmt.Sprintf("%s/%s", linkServerAddress, shortKey)
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
	return strings.Trim(URL, " ")
}
