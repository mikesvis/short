// Модуль storage для хранения в базе.
package postgres

import (
	"context"
	"database/sql"
	_goerrors "errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/keygen"
	"go.uber.org/zap"
)

type postgresDBItem struct {
	ID       string `db:"id"`
	UserID   string `db:"user_id"`
	FullURL  string `db:"full_url"`
	ShortKey string `db:"short_key"`
	Deleted  bool   `db:"is_deleted"`
}

type userUpdateItem struct {
	UserID   string
	ShortKey string
}

// Storage для хранения в базе, включает в себя указатель на sqlx.DB и логгер.
type Postgres struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

// Конструктор storage в базе. При инициализации будут созданы недостающие таблицы.
func NewPostgres(db *sqlx.DB, logger *zap.SugaredLogger) (*Postgres, error) {
	err := bootstrapDB(db)
	if err != nil {
		return nil, err
	}

	return &Postgres{db, logger}, nil
}

func bootstrapDB(db *sqlx.DB) error {
	createTableShort := `
		CREATE TABLE IF NOT EXISTS shorts (
			id varchar(36) PRIMARY KEY,
			user_id varchar(36) NOT NULL,
			full_url varchar(1000) UNIQUE NOT NULL,
			short_key varchar(255) UNIQUE NOT NULL,
			is_deleted boolean NOT NULL DEFAULT false
		)
	`
	// Как проверить эту строку? :(
	_, err := db.Exec(createTableShort)
	return err
}

// Сохранение короткой ссылки. При сохранении происходит поиск на предмет уже существующей ссылки.
// В случае если такая ссылка уже была ранее создана вернется ошибка.
func (s *Postgres) Store(ctx context.Context, u domain.URL) (domain.URL, error) {
	emptyResult := domain.URL{}

	// Как проверить эту строку? :(
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO shorts (id, user_id, full_url, short_key) VALUES ($1, $2, $3, $4) ON CONFLICT (short_key) DO NOTHING`)
	if err != nil {
		s.logger.Errorw(`Cant prepare db query`, err)
		return emptyResult, err
	}
	defer stmt.Close()

	// генерируем новый короткий урл
	item := postgresDBItem{
		ID:       uuid.NewString(),
		UserID:   u.UserID,
		FullURL:  u.Full,
		ShortKey: u.Short,
	}

	// Как проверить эту строку на ошибку? При тестах пересечение по ключу не выскакивает, хотя и не сохраняет :(
	_, err = stmt.ExecContext(ctx, item.ID, item.UserID, item.FullURL, item.ShortKey)
	if err != nil {
		s.logger.Infof("old %v", err)
		var pgErr *pgconn.PgError
		if _goerrors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
			// Ошибка непонятная
			s.logger.Errorw(`Error occured during insert`, err)
			return emptyResult, err
		}
	}

	// Ошибок не было, значит успешно сохранили с новым коротким урлом
	if err == nil {
		return u, nil
	}

	// Был конфликт пересечения по короткому урлу, забираем старый короткий урл который уже был в базе
	old, err := s.GetByFull(ctx, u.Full)
	if err != nil {
		s.logger.Errorw(`Error occured during select`, err)
		return emptyResult, err
	}

	return old, errors.ErrConflict
}

// Поиск по полной ссылке.
func (s *Postgres) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	emptyResult := domain.URL{}

	// пробуем получить по полному урлу
	// Как проверить эту строку? :(
	stmt, err := s.db.PrepareContext(ctx, `SELECT id, user_id, full_url, short_key, is_deleted FROM shorts WHERE "full_url" = $1`)
	if err != nil {
		s.logger.Errorw(`Cant prepare db query`, err)
		return emptyResult, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fullURL)

	var p postgresDBItem
	err = row.Scan(&p.ID, &p.UserID, &p.FullURL, &p.ShortKey, &p.Deleted)
	if _goerrors.Is(err, sql.ErrNoRows) {
		// нет совпадения по полному урлу, вернем пустой результат
		return emptyResult, nil
	}

	if err != nil {
		// какая-то другая ошибка
		s.logger.Errorw(`Error occured during select`, err)
		return emptyResult, err
	}

	return domain.URL{UserID: p.UserID, Full: p.FullURL, Short: p.ShortKey, Deleted: p.Deleted}, nil
}

// Поиск по короткой ссылке.
func (s *Postgres) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	emptyResult := domain.URL{}

	// пробуем получить по короткому урлу
	// Как проверить эту строку? :(
	stmt, err := s.db.PrepareContext(ctx, `SELECT id, user_id, full_url, short_key, is_deleted FROM shorts WHERE "short_key" = $1`)
	if err != nil {
		s.logger.Errorw(`Cant prepare db query`, err)
		return emptyResult, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, shortURL)

	var p postgresDBItem
	err = row.Scan(&p.ID, &p.UserID, &p.FullURL, &p.ShortKey, &p.Deleted)
	if _goerrors.Is(err, sql.ErrNoRows) {
		// нет совпадения по короткому урлу, вернем пустой результат
		return emptyResult, nil
	}

	if err != nil {
		// какая-то другая ошибка
		s.logger.Errorw(`Error occured during select`, err)
		return emptyResult, err
	}

	return domain.URL{UserID: p.UserID, Full: p.FullURL, Short: p.ShortKey, Deleted: p.Deleted}, nil
}

// Пинг базы.
func (s *Postgres) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Пакетное сохранение коротких URL. В методе используется поиск уже существующих URL.
func (s *Postgres) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
	// в мапере хранится полный урл = ключ корреляции
	mapper := make(map[string]string, len(us))
	// это хотим сохранить, но существующие будут удаляться из добавления в базу
	toStore := make(map[string]domain.URL, len(us))
	// слайс для составления select
	fullUrls := []string{}

	for k, v := range us {
		mapper[string(v.Full)] = k
		toStore[k] = v
		fullUrls = append(fullUrls, v.Full)
	}

	// какое-то неведомое колдунство? Иначе where in не сделать
	// как протестить err?
	query, args, err := sqlx.In("SELECT id, user_id, full_url, short_key, is_deleted FROM shorts WHERE full_url IN (?)", fullUrls)
	if err != nil {
		s.logger.Errorw(`Error occured while composing query`, err)
		return nil, err
	}
	query = s.db.Rebind(query)

	// ищем существующие
	// как протестить err?
	existingItems := []postgresDBItem{}
	err = s.db.SelectContext(ctx, &existingItems, query, args...)
	if err != nil {
		s.logger.Errorw(`Error occured while select`, err)
		return nil, err
	}

	for _, v := range existingItems {
		// удаляем то что сохранять не нужно
		delete(toStore, mapper[string(v.FullURL)])
		// воскрешаем старые урлы сразу в результативную мапу
		us[mapper[string(v.FullURL)]] = domain.URL{
			UserID:  v.UserID,
			Full:    v.FullURL,
			Short:   v.ShortKey,
			Deleted: v.Deleted,
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
			UserID:   v.UserID,
			FullURL:  v.Full,
			ShortKey: v.Short,
			Deleted:  v.Deleted,
		})
	}

	// сделаем добавление через транзакцию
	tx := s.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()
	// как протестить err?
	_, err = tx.NamedExecContext(ctx, `INSERT INTO shorts (id, user_id, full_url, short_key) VALUES (:id, :user_id, :full_url, :short_key)`, newItems)
	if err != nil {
		s.logger.Errorw(`Error occured while batch insert`, err)
		return nil, err
	}

	err = tx.Commit()
	// как протестить err?
	if err != nil {
		s.logger.Errorw(`Error occured while commiting transaction`, err)
		return nil, err
	}

	return us, nil
}

// Закрытие соединения к базе.
func (s *Postgres) Close() error {
	return s.db.Close()
}

// Получение ссылок, созданных пользоваетелем.
func (s *Postgres) GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	if len(userID) == 0 {
		return nil, nil
	}

	rows, err := s.db.QueryxContext(ctx, "SELECT id, user_id, full_url, short_key, is_deleted FROM shorts WHERE user_id = $1", userID)
	if err != nil {
		s.logger.Errorw(`Error occured while preparing query`, err)
		return nil, err
	}
	defer rows.Close()

	return s.fetchUserURLs(rows)
}

// не вижу особого смысла так разбивать, но раз был комментарий то вынесу вот это
func (s *Postgres) fetchUserURLs(rows *sqlx.Rows) ([]domain.URL, error) {
	p := postgresDBItem{}
	result := make([]domain.URL, 0, 20)
	for rows.Next() {
		err := rows.StructScan(&p)
		if err != nil {
			s.logger.Errorw(`Error occured while scanning row`, err)
			return nil, err
		}

		result = append(result, domain.URL{
			UserID:  p.UserID,
			Full:    p.FullURL,
			Short:   p.ShortKey,
			Deleted: p.Deleted,
		})
	}

	err := rows.Err()
	if err != nil {
		s.logger.Errorw(`Error caused by rows fetch`, err)
		return nil, err
	}

	return result, nil
}

// Пакетное удаление коротких ссылок
func (s *Postgres) DeleteBatch(ctx context.Context, userID string, pack []string) {
	inputCh := s.generator(ctx, userID, pack)
	channels := s.fanOut(ctx, inputCh)
	resultCh := s.fanIn(ctx, channels...)
	s.deleteBatchByIds(ctx, resultCh)
}

// генератор добавляет в канал сообщения
// в каждом сообщении userID + shortKey
func (s *Postgres) generator(ctx context.Context, userID string, input []string) chan userUpdateItem {
	inputCh := make(chan userUpdateItem, 10)

	go func() {
		defer close(inputCh)
		for _, data := range input {
			item := userUpdateItem{
				UserID:   userID,
				ShortKey: data,
			}
			select {
			case <-ctx.Done():
				return
			case inputCh <- item:
			}
		}
	}()

	return inputCh
}

// fanOut распределяет сообщения (userID + shortKey) на воркеры
// пишем результат воркеров в пул каналов
func (s *Postgres) fanOut(ctx context.Context, inputCh chan userUpdateItem) []chan string {
	numWorkers := 10
	channels := make([]chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		validateResCh := s.validate(ctx, inputCh)
		channels[i] = validateResCh
	}

	return channels
}

// Валидируем пользователя и не было ли уже удалено ранее
func (s *Postgres) validate(ctx context.Context, inputCh <-chan userUpdateItem) chan string {
	validateRes := make(chan string, 1)
	go func() {
		defer close(validateRes)

		for data := range inputCh {

			row := s.db.QueryRowContext(ctx, `SELECT id FROM shorts WHERE "user_id" = $1 AND "short_key" = $2 AND "is_deleted" = false`, data.UserID, data.ShortKey)
			var id string
			err := row.Scan(&id)
			if err != nil {
				s.logger.Errorw(`Error occured while scanning row`, err, `data`, data)
				continue
			}

			select {
			case <-ctx.Done():
				return
			case validateRes <- id:
			}
		}
	}()

	return validateRes
}

// fanIn объединяем каналы в результирующий канал
// в сообщениях уже только те ID которые можно update
func (s *Postgres) fanIn(ctx context.Context, resultChs ...chan string) chan string {
	finalCh := make(chan string, 10)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}

// берем из канала список id на удаление и выполняем бач update
func (s *Postgres) deleteBatchByIds(ctx context.Context, inputCh chan string) {
	var idsToDelete []string
	for idToDelete := range inputCh {
		idsToDelete = append(idsToDelete, idToDelete)
	}

	if len(idsToDelete) == 0 {
		return
	}

	// какое-то неведомое колдунство? Иначе where in не сделать
	query, args, err := sqlx.In(`UPDATE shorts SET "is_deleted" = true WHERE id IN (?)`, idsToDelete)
	if err != nil {
		s.logger.Errorw(`Error occured while making updating query`, err, `idsToDelete`, idsToDelete)
		return
	}
	query = s.db.Rebind(query)

	_, err = s.db.ExecContext(ctx, query, args...)

	if err != nil {
		s.logger.Errorw(`Error occured while updating rows`, err, `query`, query)
		return
	}
}

// Получение рандомного ключа
func (s *Postgres) GetRandkey(n uint) string {
	return keygen.GetRandkey(n)
}
