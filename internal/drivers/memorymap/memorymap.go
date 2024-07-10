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
