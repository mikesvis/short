package postgres

import (
	"context"
	"database/sql"
	_goerrors "errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
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
	// мда, чет я это... тут точно транзакция не нужна была :)
	createQuery := `
		CREATE TABLE IF NOT EXISTS short (
			id varchar(36) PRIMARY KEY,
			full_url varchar(1000) UNIQUE NOT NULL,
			short_key varchar(255) UNIQUE NOT NULL
		)
	`

	_, error := db.Exec(createQuery)

	return error
}

func (s *Postgres) Store(ctx context.Context, u domain.URL) (domain.URL, error) {
	// транзакция тут тоже не нужна
	// если получилось записать без конфликтов - то ок
	// если был конфликт, то и так не записали - нечего откатывать
	emptyResult := domain.URL{}

	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO short (id, full_url, short_key) VALUES ($1, $2, $3) ON CONFLICT (short_key) DO NOTHING`)
	if err != nil {
		return emptyResult, err
	}
	defer stmt.Close()

	// генерируем новый короткий урл
	item := postgresDBItem{
		ID:       uuid.NewString(),
		FullURL:  u.Full,
		ShortKey: u.Short,
	}

	_, err = stmt.ExecContext(ctx, item.ID, item.FullURL, item.ShortKey)
	if err != nil {
		var pgErr *pgconn.PgError
		if _goerrors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
			// Ошибка непонятная
			return emptyResult, err
		}
	}

	// Ошибок не было, значит успешно сохранили с новым коротким урлом
	if err == nil {
		return u, nil
	}

	// Была ошибка пересечения по короткому урлу, забираем старый короткий урл который уже был в базе
	// Как я не пробовал избавится от этого селекта - не смог
	// RETURNING id - работает только если вставили без ошибок
	// UPSERT делать нельзя - поскольку мне нечего исключать из SET (нельзя перезаписывать старый короткий урл,
	// а зачем перезаписывать id для чего? Все равно придется делать SELECT по нему для получение старого короткого урла
	old, err := s.GetByFull(ctx, u.Full)
	if err != nil {
		return emptyResult, err
	}

	return old, errors.ErrConflict
}

func (s *Postgres) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	emptyResult := domain.URL{}

	// пробуем получить по полному урлу
	stmt, err := s.db.PrepareContext(ctx, `SELECT id, full_url, short_key FROM short WHERE "full_url" = $1`)
	if err != nil {
		return emptyResult, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fullURL)

	var p postgresDBItem
	err = row.Scan(&p.ID, &p.FullURL, &p.ShortKey)
	if _goerrors.Is(err, sql.ErrNoRows) {
		// нет совпадения по полному урлу, вернем пустой результат
		return emptyResult, nil
	}

	if err != nil {
		// какая-то другая ошибка
		return emptyResult, err
	}

	return domain.URL{Full: p.FullURL, Short: p.ShortKey}, nil
}

func (s *Postgres) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	emptyResult := domain.URL{}

	// пробуем получить по короткому урлу
	stmt, err := s.db.PrepareContext(ctx, `SELECT id, full_url, short_key FROM short WHERE "short_key" = $1`)
	if err != nil {
		return emptyResult, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, shortURL)

	var p postgresDBItem
	err = row.Scan(&p.ID, &p.FullURL, &p.ShortKey)
	if _goerrors.Is(err, sql.ErrNoRows) {
		// нет совпадения по короткому урлу, вернем пустой результат
		return emptyResult, nil
	}

	if err != nil {
		// какая-то другая ошибка
		return emptyResult, err
	}

	return domain.URL{Full: p.FullURL, Short: p.ShortKey}, nil
}

func (s *Postgres) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
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

func (s *Postgres) Close() error {
	return s.db.Close()
}
