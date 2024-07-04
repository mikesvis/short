package memorymap

import (
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

func (s *MemoryMap) Store(u domain.URL) error {
	s.items[domain.ID(uuid.NewString())] = u
	return nil
}

func (s *MemoryMap) GetByFull(fullURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Full) != fullURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}

func (s *MemoryMap) GetByShort(shortURL string) (domain.URL, error) {
	for _, v := range s.items {
		if string(v.Short) != shortURL {
			continue
		}

		return v, nil
	}

	return domain.URL{}, nil
}
