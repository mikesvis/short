package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/drivers/filedb"
	"github.com/mikesvis/short/internal/drivers/memorymap"
	"github.com/mikesvis/short/internal/drivers/postgres"
)

type Storage interface {
	Store(context.Context, domain.URL) error
	GetByFull(ctx context.Context, fullURL string) (domain.URL, error)
	GetByShort(ctx context.Context, shortURL string) (domain.URL, error)
	Ping(ctx context.Context) error
}

func NewStorage(c *config.Config) Storage {
	if len(string(c.DatabaseDSN)) != 0 {
		db, err := sql.Open("pgx", string(c.DatabaseDSN))
		if err != nil {
			panic(err)
		}
		// defer db.Close() - вот с закрытием немного не ясно, поскольку тут закрытию точно не место
		// а делать для каждого типа storage Close - странно, т.к. не все типы требуют закрытия. Пожоже на протекание абстракции.
		// Кто-то в инете говорит что закрывать соединение не надо (гуглил на stackoverflow), мол
		// при завершении приложения все коннекты и так закрываются и переживать не надо (в отличии от rows.close()).
		// Если соединения закрывать все-таки нужно, то я бы сделал некую структуру с списком соединений.
		// Данная структура наполнялась бы соединениями (db / http / etc) в момент инициализации в app или main
		// и в defer верхнего уровня закрывала бы соединения из списка.

		return postgres.NewPostgres(db)
	}

	if len(string(c.FileStoragePath)) != 0 {
		return filedb.NewFileDB(string(c.FileStoragePath))
	}

	return memorymap.NewMemoryMap()
}
