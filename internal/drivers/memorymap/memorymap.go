package memorymap

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
)

type MemoryMap struct {
	items map[domain.ID]domain.URL
}

func NewMemoryMap() *MemoryMap {
	items := make(map[domain.ID]domain.URL)
	return &MemoryMap{items: items}
}

func (s *MemoryMap) Store(ctx context.Context, u domain.URL) error {
	s.items[domain.ID(uuid.NewString())] = u
	return nil
}

func (s *MemoryMap) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Full) != fullURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}

func (s *MemoryMap) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Short) != shortURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}

func (s *MemoryMap) Ping(ctx context.Context) error {
	return nil
}

func (s *MemoryMap) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
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
				Full:  v.Full,
				Short: v.Short,
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
