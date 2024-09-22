// Модуль хранилки.
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

// Интерфейс хранилки
type Storage interface {
	// Сохранение короткой ссылки.
	Store(ctx context.Context, URL domain.URL) (domain.URL, error)

	// Пакетное сохранение коротких ссылок.
	StoreBatch(ctx context.Context, pack map[string]domain.URL) (map[string]domain.URL, error)

	// Получение короткой ссылки по полной.
	GetByFull(ctx context.Context, fullURL string) (domain.URL, error)

	// Получение полной ссылки по короткой.
	GetByShort(ctx context.Context, shortURL string) (domain.URL, error)

	// Получение ссылок пользователя.
	GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error)
}

// Интерфейс обеспечивающий метод для прозвона хранилки.
type StoragePinger interface {
	Storage
	// Прозвон хранилки.
	Ping(ctx context.Context) error
}

// Интерфейс обеспечивающий метод для закрытия хранилки.
type StorageCloser interface {
	Storage
	// Закрытие хранилки.
	Close() error
}

// Интерфейс обеспечивающий метод для удаления URL из хранилки.
type StorageDeleter interface {
	Storage
	// Пакетное удаление URL.
	DeleteBatch(ctx context.Context, userID string, pack []string)
}

// Интерфейс, объединяющий прозвон, закрытие и пакетное удаление.
type StoragePingerCloserDeleter interface {
	StoragePinger
	StorageCloser
	StorageDeleter
}

// Конструктор хранилки. На основании конфига выбирается движок для хранилки.
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
