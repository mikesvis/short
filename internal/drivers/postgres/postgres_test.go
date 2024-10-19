// Модуль storage для хранения в базе.
package postgres

import (
	_context "context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
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

func TestPostgres_Store(t *testing.T) {
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	defer db.Close()

	type args struct {
		u    domain.URL
		aCtx _context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    domain.URL
		wantErr bool
	}{
		{
			name: "Store item",
			args: args{
				u: domain.URL{
					UserID: "DoomGuy",
					Full:   "http://www.yandex.ru/verylongpath",
					Short:  "short",
				},
				aCtx: _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
			},
			want: domain.URL{
				UserID:  "DoomGuy",
				Full:    "http://www.yandex.ru/verylongpath",
				Short:   "short",
				Deleted: false,
			},
			wantErr: false,
		},
		{
			name: "Has conflict in store",
			args: args{
				u: domain.URL{
					UserID: "Heretic",
					Full:   "http://www.yandex.ru/verylongpath",
					Short:  "short",
				},
				aCtx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Heretic"),
			},
			want:    domain.URL{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewPostgres(db, l)
			if tt.wantErr {
				db.Close()
			}
			got, err := s.Store(tt.args.aCtx, tt.args.u)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
