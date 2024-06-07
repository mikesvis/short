package storage

import (
	"github.com/mikesvis/short/internal/domain"
)

type StorageURL interface {
	Store(domain.URL)
	GetByFull(fullURL string) domain.URL
	GetByShort(shortURL string) domain.URL
}

type storageURL struct {
	items map[string]domain.URL
}

func NewStorageURL(items map[string]domain.URL) StorageURL {
	return &storageURL{items: items}
}

func (s *storageURL) Store(u domain.URL) {
	s.items[string(u.Short)] = u
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
