// Модуль форматирования ссылок
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
func ValidateURL(urlToValidate string) error {
	if len(urlToValidate) == 0 {
		return fmt.Errorf("URL can not be empty")
	}

	_, err := url.ParseRequestURI(urlToValidate)
	if err != nil {
		return fmt.Errorf("URL is not an URL format, %s given", err)
	}

	return nil
}

// Чистка URL
func SanitizeURL(urlToTrim string) string {
	return strings.Trim(urlToTrim, " ")
}
