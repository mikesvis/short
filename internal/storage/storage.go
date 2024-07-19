package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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
}

type StoragePinger interface {
	Storage
	Ping(ctx context.Context) error
}

type StorageCloser interface {
	Storage
	Close() error
}

type StoragePingerCloser interface {
	StoragePinger
	StorageCloser
}

func NewStorage(c *config.Config) Storage {
	if len(string(c.DatabaseDSN)) != 0 {
		db, err := sqlx.Open("postgres", string(c.DatabaseDSN))
		if err != nil {
			panic(err)
		}

		return StoragePingerCloser(postgres.NewPostgres(db))
	}

	if len(string(c.FileStoragePath)) != 0 {
		return StoragePinger(filedb.NewFileDB(string(c.FileStoragePath)))
	}

	return inmemory.NewInMemory()
}
