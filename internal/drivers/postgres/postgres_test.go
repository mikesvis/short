// Модуль storage для хранения в базе.
package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikesvis/short/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getDataBaseDSN() string {
	return "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
}

func TestNewPostgres(t *testing.T) {
	l, _ := logger.NewLogger()
	type args struct {
		databaseDSN string
	}
	tests := []struct {
		name string
		args args
		want *Postgres
	}{
		{
			name: "New storage is of type",
			args: args{
				databaseDSN: getDataBaseDSN(),
			},
			want: &Postgres{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqlx.Open("postgres", tt.args.databaseDSN)
			require.NoError(t, err)
			newStorage, err := NewPostgres(db, l)
			require.NoError(t, err)
			assert.IsType(t, tt.want, newStorage)
		})
	}
}
