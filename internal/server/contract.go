package server

import "github.com/mikesvis/short/internal/domain"

type StorageURL interface {
	Store(domain.URL)
	GetByFull(fullURL string) domain.URL
	GetByShort(shortURL string) domain.URL
}
