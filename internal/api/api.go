package api

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/mikesvis/short/pkg/urlformat"
)

type URL string

type Request struct {
	URL URL `json:"url"`
}

type Response struct {
	Result URL `json:"result"`
}

type CorrelationID string
type OriginalURL string
type ShortURL string

type BatchRequest struct {
	CorrelationID CorrelationID `json:"correlation_id"`
	OriginalURL   OriginalURL   `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID CorrelationID `json:"correlation_id"`
	ShortURL      ShortURL      `json:"short_url"`
}

// А тут ли это должно быть?
func (c *BatchRequest) UnmarshalJSON(data []byte) error {
	type BatchRequestAlias BatchRequest

	tmp := &struct {
		*BatchRequestAlias
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}{
		BatchRequestAlias: (*BatchRequestAlias)(c),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	tmp.CorrelationID = strings.TrimSpace(tmp.CorrelationID)

	if len(string(tmp.CorrelationID)) == 0 {
		return errors.New(`correlation_id can not be empty`)
	}

	tmp.OriginalURL = strings.TrimSpace(tmp.OriginalURL)

	if len(string(tmp.OriginalURL)) == 0 {
		return errors.New(`original_url can not be empty`)
	}

	if err := urlformat.ValidateURL(tmp.OriginalURL); err != nil {
		return err
	}

	c.CorrelationID = CorrelationID(tmp.CorrelationID)
	c.OriginalURL = OriginalURL(tmp.OriginalURL)

	return nil
}
