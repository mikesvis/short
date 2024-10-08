// Модуль api служит для описания структур запросов/ответов
package api

// URL - полный адрес URL в виде строки
type URL string

// Request - запрос с полем URL, которое требуется сократить в JSON формате
type Request struct {
	URL URL `json:"url"`
}

// Resonse - ответ в JSON формате с коротким URL
type Response struct {
	Result URL `json:"result"`
}

// BatchRequest - запрос с пакетным сокращением URL
type BatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse - ответ с пакетным сокращением URL
type BatchResponse []struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserResponse - ответ с сокращенными и изначальными URL пользователя
type UserResponse []struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// BatchDeleteRequest - запрос на пакетное удаление скоращенных URL
type BatchDeleteRequest []string
