package server

import "github.com/mikesvis/short/internal/domain"

type StorageURL interface {
	Store(domain.URL)
	GetByFull(fullURL string) (domain.URL, error)
	GetByShort(shortURL string) (domain.URL, error)
}
