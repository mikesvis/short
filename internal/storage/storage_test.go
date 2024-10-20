// Модуль хранилки.
package storage

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/drivers/filedb"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/drivers/postgres"
	"github.com/mikesvis/short/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Вот эта хрень потому что я не допетрил как победить автотесты
// Поэтому я у себя ставлю хост 0.0.0.0 а для долбанного github postgres
// Если знаешь как решить - скажи
// Дело в то что при автотестах поднимается контейнер и к нему можно подключиться по
// postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable
// Как сделать также без изменения файла hosts я хз, поэтому и проблема
func getDataBaseDSN() string {
	// return "postgres://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable"
	return "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
}

func TestNewStorage(t *testing.T) {
	logger, _ := logger.NewLogger()
	type args struct {
		c *config.Config
	}

	type want struct {
		wantErr bool
		Storage
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "In memory storage",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "",
				},
			},
			want: want{
				wantErr: false,
				Storage: &inmemory.InMemory{},
			},
		},
		{
			name: "In file storage",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "/tmp/db.tmp",
					DatabaseDSN:     "",
				},
			},
			want: want{
				wantErr: false,
				Storage: &filedb.FileDB{},
			},
		},
		{
			name: "Postgres storage",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     getDataBaseDSN(),
				},
			},
			want: want{
				wantErr: false,
				Storage: &postgres.Postgres{},
			},
		},
		{
			name: "Postgres storage failed",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "dummyDb",
				},
			},
			want: want{
				wantErr: true,
				Storage: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStorage(tt.args.c, logger)
			if tt.want.wantErr {
				require.Error(t, err)
				return

			}

			require.NoError(t, err)
			assert.IsType(t, tt.want.Storage, s)
		})
	}
}
