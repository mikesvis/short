package api

type URL string

type Request struct {
	URL URL `json:"url"`
}

type Response struct {
	Result URL `json:"result"`
}

type BatchRequest []BatchRequestItem

type BatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse []BatchResponseItem

type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserResponse []UserResponseItem

type UserResponseItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchDeleteRequest []string
