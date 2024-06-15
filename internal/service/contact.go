package service

import "github.com/mikesvis/short/internal/domain"

type storageURL interface {
	Store(domain.URL)
	GetByFull(fullURL string) domain.URL
	GetByShort(shortURL string) domain.URL
}
