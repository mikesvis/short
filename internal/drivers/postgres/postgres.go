package postgres

import (
	"context"
	"database/sql"
	_goerrors "errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
)

type postgresDBItem struct {
	ID       string `db:"id"`
	FullURL  string `db:"full_url"`
	ShortKey string `db:"short_key"`
}

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	err := bootstrapDB(db)
	if err != nil {
		panic(err)
	}

	return &Postgres{db}
}

func bootstrapDB(db *sqlx.DB) error {
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
	// если было конфликта, то и так не записали - нечего откатывать
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

	// Был конфликт пересечения по короткому урлу, забираем старый короткий урл который уже был в базе
	// Как я не пробовал избавится от этого селекта - не смог
	// RETURNING id - работает только если вставили без ошибок
	// UPSERT делать нельзя - поскольку мне нечего исключать из SET (нельзя перезаписывать старый короткий урл,
	// а зачем перезаписывать id для чего? Все равно придется делать SELECT по нему для получение старого короткого урла
	// эта рализация была для предыдущего драйвера, в StoreBatch использую  sqlx
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
	// в мапере хранится полный урл = ключ корреляции
	mapper := make(map[string]string, len(us))
	// это хотим сохранить, но существующие будут удаляться из добавления в базу
	toStore := make(map[string]domain.URL, len(us))
	// слайс для составления select
	fullUrls := []string{}

	for k, v := range us {
		mapper[string(v.Full)] = k
		fullUrls = append(fullUrls, v.Full)
		toStore[k] = v
	}

	// какое-то неведомое колдунство? Иначе where in не сделать
	query, args, err := sqlx.In("SELECT id, full_url, short_key FROM short WHERE full_url IN (?)", fullUrls)
	if err != nil {
		return nil, err
	}
	query = s.db.Rebind(query)

	// ищем существующие
	existingItems := []postgresDBItem{}
	err = s.db.SelectContext(ctx, &existingItems, query, args...)
	if err != nil {
		return nil, err
	}

	for _, v := range existingItems {
		// удаляем то что сохранять не нужно
		delete(toStore, mapper[string(v.FullURL)])
		// воскрешаем старые урлы сразу в результативную мапу
		us[mapper[string(v.FullURL)]] = domain.URL{
			Full:  v.FullURL,
			Short: v.ShortKey,
		}
	}

	// нечего сохранять - уходим
	if len(toStore) == 0 {
		return us, nil
	}

	// запоняем структуры для сохранения новых данных
	newItems := []postgresDBItem{}
	for _, v := range toStore {
		newItems = append(newItems, postgresDBItem{
			ID:       uuid.NewString(),
			FullURL:  v.Full,
			ShortKey: v.Short,
		})
	}

	// сделаем добавление через транзакцию
	tx := s.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx, `INSERT INTO short (id, full_url, short_key) VALUES (:id, :full_url, :short_key)`, newItems)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (s *Postgres) Close() error {
	return s.db.Close()
}
