package storage

import (
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/filedb"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/memorymap"
)

type Storage interface {
	Store(domain.URL)
	GetByFull(fullURL string) (domain.URL, error)
	GetByShort(shortURL string) (domain.URL, error)
}

func NewStorage(fileStoragePath string) Storage {
	if len(fileStoragePath) == 0 {
		logger.Log.Info("Using in-memory map storage")
		return memorymap.NewMemoryMap()
	}

	logger.Log.Infof("Using file storage by path %s", fileStoragePath)
	return filedb.NewFileDb(fileStoragePath)
}
