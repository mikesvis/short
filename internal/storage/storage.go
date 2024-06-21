package storage

import (
	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
)

type storageURL struct {
	items map[domain.ID]domain.URL
}

func NewStorageURL(items map[domain.ID]domain.URL) *storageURL {
	return &storageURL{items: items}
}

func (s *storageURL) Store(u domain.URL) {
	s.items[domain.ID(uuid.NewString())] = u
}

func (s *storageURL) GetByFull(fullURL string) domain.URL {
	for _, v := range s.items {
		if string(v.Full) != fullURL {
			continue
		}

		return v
	}

	return domain.URL{}
}

func (s *storageURL) GetByShort(shortURL string) domain.URL {
	for _, v := range s.items {
		if string(v.Short) != shortURL {
			continue
		}

		return v
	}

	return domain.URL{}
}
