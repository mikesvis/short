package storage

import (
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/filedb"
	"github.com/mikesvis/short/internal/memorymap"
)

type Storage interface {
	Store(domain.URL) error
	GetByFull(fullURL string) (domain.URL, error)
	GetByShort(shortURL string) (domain.URL, error)
}

func NewStorage(fileStoragePath string) Storage {
	if len(fileStoragePath) == 0 {
		return memorymap.NewMemoryMap()
	}

	return filedb.NewFileDB(fileStoragePath)
}
