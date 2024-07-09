package postgres

import (
	"database/sql"

	"github.com/mikesvis/short/internal/domain"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db}
}

func (s *Postgres) Store(u domain.URL) error {
	return nil
}

func (s *Postgres) GetByFull(fullURL string) (domain.URL, error) {
	return domain.URL{}, nil
}

func (s *Postgres) GetByShort(shortURL string) (domain.URL, error) {
	return domain.URL{}, nil
}

func (s *Postgres) Ping() error {
	return s.db.Ping()
}
