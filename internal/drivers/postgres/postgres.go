package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
)

type postgresDBItem struct {
	ID       string
	FullURL  string
	ShortKey string
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	err := bootstrapDB(db)
	if err != nil {
		panic(err)
	}

	return &Postgres{db}
}

func bootstrapDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	createQuery := `
		CREATE TABLE IF NOT EXISTS short (
			id varchar(36) PRIMARY KEY,
			full_url varchar(1000) UNIQUE NOT NULL,
			short_key varchar(255) UNIQUE NOT NULL
		)
	`
	tx.Exec(createQuery)

	return tx.Commit()
}

func (s *Postgres) Store(ctx context.Context, u domain.URL) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO short (id, full_url, short_key) VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	item := postgresDBItem{
		ID:       uuid.NewString(),
		FullURL:  u.Full,
		ShortKey: u.Short,
	}

	_, err = stmt.ExecContext(ctx, item.ID, item.FullURL, item.ShortKey)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Postgres) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	return s.findByColumn(ctx, "full_url", fullURL)
}

func (s *Postgres) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	return s.findByColumn(ctx, "short_key", shortURL)
}

func (s *Postgres) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Postgres) findByColumn(ctx context.Context, column, needle string) (domain.URL, error) {
	// вот это какой-то бред! как я не пытался передать в формирование запроса prepare колонку, не смог
	var countQuery, selectQuery string
	switch column {
	case "full_url":
		countQuery = `SELECT COUNT(1) FROM short WHERE "full_url" = $1`
		selectQuery = `SELECT id, full_url, short_key FROM short WHERE "full_url" = $1`
	case "short_key":
		countQuery = `SELECT COUNT(1) FROM short WHERE "short_key" = $1`
		selectQuery = `SELECT id, full_url, short_key FROM short WHERE "short_key" = $1`
	}

	stmt, err := s.db.PrepareContext(ctx, countQuery)
	if err != nil {
		return domain.URL{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, needle)

	var rowsNum int

	err = row.Scan(&rowsNum)
	if err != nil {
		return domain.URL{}, err
	}

	if rowsNum == 0 {
		return domain.URL{}, nil
	}

	stmt, err = s.db.PrepareContext(ctx, selectQuery)
	if err != nil {
		return domain.URL{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowContext(ctx, needle)

	var p postgresDBItem

	err = row.Scan(&p.ID, &p.FullURL, &p.ShortKey)
	if err != nil {
		return domain.URL{}, err
	}

	return domain.URL{Full: p.FullURL, Short: p.ShortKey}, nil
}

func (s *Postgres) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
	// подход в рамках одной транзации: берем 1 элемент, смотрим есть ли он в базе,
	// если есть, то ничего не делаем, просто замеяем его урл на старый
	// если нет - пишем в базу
	// мое мнение: это убогое решение, поскольку будет порождено столько транзакций/select'ов/insert'ов
	// сколько элементов пришло в метод, реплике может стать очень плохо в перспективе кол-ва эл-тов на входе
	// надо делать завдержку реплики
	for k, v := range us {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		selectQuery := `SELECT id, full_url, short_key FROM short WHERE "full_url" = $1`
		stmt, err := tx.PrepareContext(ctx, selectQuery)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		row := stmt.QueryRowContext(ctx, v.Full)

		var p postgresDBItem
		err = row.Scan(&p.ID, &p.FullURL, &p.ShortKey)
		if err != sql.ErrNoRows && err != nil {
			return nil, err
		}

		// такой элемент уже есть, не добавляем, узнаем его старый урл
		if err == nil {
			us[k] = domain.URL{
				Full:  v.Full,
				Short: p.ShortKey,
			}
			continue
		}

		// этот элемент новый, будем его сохранять
		stmt, err = tx.PrepareContext(ctx, `INSERT INTO short (id, full_url, short_key) VALUES ($1, $2, $3)`)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		item := postgresDBItem{
			ID:       uuid.NewString(),
			FullURL:  v.Full,
			ShortKey: v.Short,
		}

		_, err = stmt.ExecContext(ctx, item.ID, item.FullURL, item.ShortKey)
		if err != nil {
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return us, nil
}
