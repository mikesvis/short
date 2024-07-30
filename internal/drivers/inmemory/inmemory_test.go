package inmemory

import (
	_context "context"
	"math/rand"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewStorageURL(t *testing.T) {
	type args struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name string
		args args
		want *InMemory
	}{
		{
			name: "New storage is of type",
			args: args{
				items: map[domain.ID]domain.URL{},
			},
			want: &InMemory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStorage := NewInMemory()
			assert.IsType(t, tt.want, newStorage)
		})
	}
}

func Test_storageURL_Store(t *testing.T) {
	ctx := _context.Background()

	type fields struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name   string
		fields fields
		args   domain.URL
		want   map[domain.ID]domain.URL
	}{
		{
			name: "Store item",
			fields: fields{
				items: map[domain.ID]domain.URL{},
			},
			args: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "http://localhost/short",
			},
			want: map[domain.ID]domain.URL{
				"52fdfc07-2182-454f-963f-5f0f9a621d72": {
					Full:    "http://www.yandex.ru/verylongpath",
					Short:   "http://localhost/short",
					Deleted: false,
				},
			},
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemory{
				items: tt.fields.items,
			}
			s.Store(ctx, tt.args)
			assert.EqualValues(t, tt.want, s.items)
		})
	}
}

func Test_storageURL_GetByFull(t *testing.T) {
	ctx := _context.Background()

	type fields struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   domain.URL
	}{
		{
			name: "Is found by full",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath1",
						Short: "http://localhost/short1",
					},
					"dummyId2": {
						Full:  "http://www.google.ru/verylongpath2",
						Short: "http://localhost/short2",
					},
				},
			},
			args: "http://www.yandex.ru/verylongpath1",
			want: domain.URL{
				Full:    "http://www.yandex.ru/verylongpath1",
				Short:   "http://localhost/short1",
				Deleted: false,
			},
		}, {
			name: "Is not found by full",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath1",
						Short: "http://localhost/short1",
					},
				},
			},
			args: "http://localhost/",
			want: domain.URL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemory{
				items: tt.fields.items,
			}
			item, _ := s.GetByFull(ctx, tt.args)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
		})
	}
}

func Test_storageURL_GetByShort(t *testing.T) {
	ctx := _context.Background()

	type fields struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   domain.URL
	}{
		{
			name: "Is found by short",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath1",
						Short: "http://localhost/short1",
					},
					"dummyId2": {
						Full:  "http://www.google.ru/verylongpath2",
						Short: "http://localhost/short2",
					},
				},
			},
			args: "http://localhost/short1",
			want: domain.URL{
				Full:    "http://www.yandex.ru/verylongpath1",
				Short:   "http://localhost/short1",
				Deleted: false,
			},
		}, {
			name: "Is not found by short",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"dummyId1": {
						Full:  "http://www.yandex.ru/verylongpath1",
						Short: "http://localhost/short1",
					},
				},
			},
			args: "http://localhost/",
			want: domain.URL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemory{
				items: tt.fields.items,
			}
			item, _ := s.GetByShort(ctx, tt.args)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
		})
	}
}

func TestMemoryMap_StoreBatch(t *testing.T) {
	ctx := _context.Background()

	type fields struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name   string
		fields fields
		args   map[string]domain.URL
		want   map[string]domain.URL
	}{
		{
			name: "Batch store items",
			fields: fields{
				items: map[domain.ID]domain.URL{},
			},
			args: map[string]domain.URL{
				"1": {
					Full:  "http://www.yandex.ru/verylongpath1",
					Short: "short1",
				},
			},
			want: map[string]domain.URL{
				"1": {
					Full:    "http://www.yandex.ru/verylongpath1",
					Short:   "short1",
					Deleted: false,
				},
			},
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemory{
				items: tt.fields.items,
			}
			stored, err := s.StoreBatch(ctx, tt.args)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, stored)
		})
	}
}

func TestInMemory_GetUserURLs(t *testing.T) {
	ctx := _context.WithValue(_context.Background(), context.ContextUserKey, "DoomGuy")
	type fields struct {
		items map[domain.ID]domain.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		want    []domain.URL
		wantErr bool
	}{
		{
			name: "Get current user URLs",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"1": {
						UserID: "DoomGuy",
						Full:   "http://iddqd.com",
						Short:  "idkfa",
					},
				},
			},
			args: "DoomGuy",
			want: []domain.URL{
				{
					UserID:  "DoomGuy",
					Full:    "http://iddqd.com",
					Short:   "idkfa",
					Deleted: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Get empty list user URLs",
			fields: fields{
				items: map[domain.ID]domain.URL{
					"1": {
						UserID: "DoomGuy",
						Full:   "http://iddqd.com",
						Short:  "idkfa",
					},
				},
			},
			args:    "Heretic",
			want:    []domain.URL{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemory{
				items: tt.fields.items,
			}
			got, err := s.GetUserURLs(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemory.GetUserURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemory.GetUserURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
