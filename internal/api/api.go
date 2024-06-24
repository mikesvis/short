package api

type URL string

type Request struct {
	URL URL `json:"url"`
}

type Response struct {
	Result URL `json:"result"`
}
