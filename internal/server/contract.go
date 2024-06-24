package server

import "github.com/mikesvis/short/internal/domain"

type StorageURL interface {
	Store(domain.URL)
	GetByFull(fullURL string) (domain.URL, bool)
	GetByShort(shortURL string) (domain.URL, bool)
}
