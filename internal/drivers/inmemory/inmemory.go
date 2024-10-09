// Модуль storage для хранения в памяти.
package inmemory

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/keygen"
	"go.uber.org/zap"
)

// Storage для хранения в памяти, включает в себя мапу с элементами ссылок и логгер.
type InMemory struct {
	items  map[domain.ID]domain.URL
	logger *zap.SugaredLogger
}

// Конструктор storage в памяти.
func NewInMemory(logger *zap.SugaredLogger) *InMemory {
	items := make(map[domain.ID]domain.URL)
	return &InMemory{items, logger}
}

// Сохранение короткой ссылки. При сохранении происходит поиск на предмет уже существующей ссылки.
// В случае если такая ссылка уже была ранее создана вернется ошибка.
func (s *InMemory) Store(ctx context.Context, u domain.URL) (domain.URL, error) {
	for _, v := range s.items {
		if v.Full == u.Full {
			return v, errors.ErrConflict
		}
	}
	s.items[domain.ID(uuid.NewString())] = u
	return u, nil
}

// Поиск по полной ссылке.
func (s *InMemory) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Full) != fullURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}

// Поиск по короткой ссылке.
func (s *InMemory) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Short) != shortURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}

// Пакетное сохранение коротких URL. В методе используется поиск уже существующих URL.
func (s *InMemory) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
	// в мапе хранится полный урл = ключ корреляции
	wantToStore := make(map[string]string, len(us))

	for k, v := range us {
		wantToStore[string(v.Full)] = k
	}

	// для начала найдем совпадения по урлу, которые были сохранены ранее
	for _, v := range s.items {
		k, exists := wantToStore[v.Full]
		if exists {
			// урл был сохранен ранее: удаляем из списка на сохранение и
			// восстанавливаем его старый short вместо нового
			delete(wantToStore, v.Full)
			us[k] = domain.URL{
				UserID:  v.UserID,
				Full:    v.Full,
				Short:   v.Short,
				Deleted: v.Deleted,
			}
		}

		// список на сохранение пустой, не смысла искать далее (все элементы уже есть в хранилке)
		if len(wantToStore) == 0 {
			break
		}
	}

	// все элементы на сохранение уже есть, нечего сохранять
	if len(wantToStore) == 0 {
		return us, nil
	}

	// будем сохранять только те елементы, которых нет
	for _, v := range wantToStore {
		s.items[domain.ID(uuid.NewString())] = us[v]
	}

	return us, nil
}

// Получение ссылок, созданных пользоваетелем.
func (s *InMemory) GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	result := make([]domain.URL, 0, 20)
	for _, v := range s.items {
		if v.UserID != userID {
			continue
		}

		result = append(result, v)
	}

	return result, nil
}

// Получение рандомного ключа
func (s *InMemory) GetRandkey(n uint) string {
	return keygen.GetRandkey(n)
}
