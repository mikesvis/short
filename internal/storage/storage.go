package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/drivers/filedb"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/drivers/postgres"
)

type Storage interface {
	Store(ctx context.Context, URL domain.URL) (domain.URL, error)
	StoreBatch(ctx context.Context, pack map[string]domain.URL) (map[string]domain.URL, error)
	GetByFull(ctx context.Context, fullURL string) (domain.URL, error)
	GetByShort(ctx context.Context, shortURL string) (domain.URL, error)
	GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error)
}

type StoragePinger interface {
	Storage
	Ping(ctx context.Context) error
}

type StorageCloser interface {
	Storage
	Close() error
}

type StorageDeleter interface {
	Storage
	DeleteBatch(ctx context.Context, userID string, pack []string)
}

type StoragePingerCloserDeleter interface {
	StoragePinger
	StorageCloser
	StorageDeleter
}

func NewStorage(c *config.Config, logger *zap.SugaredLogger) Storage {
	if len(string(c.DatabaseDSN)) != 0 {
		db, err := sqlx.Open("postgres", string(c.DatabaseDSN))
		if err != nil {
			panic(err)
		}

		return StoragePingerCloserDeleter(postgres.NewPostgres(db, logger))
	}

	if len(string(c.FileStoragePath)) != 0 {
		return StoragePinger(filedb.NewFileDB(string(c.FileStoragePath), logger))
	}

	return inmemory.NewInMemory(logger)
}
