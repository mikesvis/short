// Модуль хранилки.
package storage

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/drivers/filedb"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/drivers/postgres"
	"github.com/mikesvis/short/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewStorage(t *testing.T) {
	logger, _ := logger.NewLogger()
	dbDSN := "host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable"
	db, _ := sqlx.Open("postgres", string(dbDSN))
	postgresStorage, _ := postgres.NewPostgres(db, logger)
	type args struct {
		c      *config.Config
		logger *zap.SugaredLogger
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
				Storage: inmemory.NewInMemory(logger),
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
				Storage: filedb.NewFileDB("/tmp/db.tmp", logger),
			},
		},
		{
			name: "In file storage success",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable",
				},
			},
			want: want{
				wantErr: false,
				Storage: postgresStorage,
			},
		},
		{
			name: "In file storage fail",
			args: args{
				c: &config.Config{
					ServerAddress:   "127.0.0.1",
					BaseURL:         "http://short.go",
					FileStoragePath: "",
					DatabaseDSN:     "host=0.0.0.0 port=54344 user=postgres password=postgres dbname=short sslmode=disable",
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
			s, err := NewStorage(tt.args.c, tt.args.logger)
			if !tt.want.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			assert.ObjectsAreEqual(tt.want, s)
		})
	}
}
