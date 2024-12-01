// Модуль storage для хранения в базе.
package postgres

import (
	_context "context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/keygen"
	"github.com/mikesvis/short/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Вот эта хрень нужна потому что я не допетрил как победить автотесты на гитхаб
// Поэтому я у себя ставлю хост 0.0.0.0 а для долбанного github указан в тестах хост postgres
// Если знаешь как решить - скажи
// Дело в том что при автотестах поднимается контейнер и к нему можно подключиться по
// postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable
// Как сделать также без изменения файла hosts локально я хз, поэтому и проблема
// Пробовал указывать и networks и host/domain в docker-compose - не помогает
func getDataBaseDSN() string {
	// return "postgres://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable"
	return "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
}

func truncateTable(db *sqlx.DB) {
	db.Exec(`TRUNCATE TABLE shorts`)
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

func TestFileDB_GetRandkey(t *testing.T) {
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	defer db.Close()
	type want struct {
		typeOf  string
		len     int
		isEmpty bool
	}
	tests := []struct {
		name string
		arg  uint
		want want
	}{
		{
			name: "Rand key is string of length",
			arg:  5,
			want: want{
				typeOf:  "",
				len:     5,
				isEmpty: false,
			},
		}, {
			name: "Rand key is empty sting of zero length",
			arg:  0,
			want: want{
				typeOf:  "",
				len:     0,
				isEmpty: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     db,
				logger: l,
			}
			randKey := s.GetRandkey(tt.arg)
			assert.IsType(t, "", randKey)
			assert.Len(t, randKey, tt.want.len)
			if !tt.want.isEmpty {
				assert.NotEmpty(t, s.GetRandkey(tt.arg))
				return
			}

			assert.Empty(t, s.GetRandkey(tt.arg))
		})
	}
}

func TestPostgres_GetByFull(t *testing.T) {
	rndString1 := keygen.GetRandkey(5)
	rndString2 := keygen.GetRandkey(5)
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	defer db.Close()
	s := &Postgres{
		db:     db,
		logger: l,
	}
	s.Store(
		_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
		domain.URL{
			UserID:  "DoomGuy",
			Full:    `https://` + rndString2 + `.com`,
			Short:   rndString2,
			Deleted: false,
		},
	)
	type fields struct {
		db     *sqlx.DB
		logger *zap.SugaredLogger
	}
	type args struct {
		ctx     _context.Context
		fullURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.URL
		wantErr bool
	}{
		{
			name: "Is not found by full - empty result",
			fields: fields{
				db:     db,
				logger: l,
			},
			args: args{
				ctx:     _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				fullURL: `https://` + rndString1 + `.com`,
			},
			want:    domain.URL{},
			wantErr: false,
		},
		{
			name: "Is found by full",
			fields: fields{
				db:     db,
				logger: l,
			},
			args: args{
				ctx:     _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				fullURL: `https://` + rndString2 + `.com`,
			},
			want: domain.URL{
				UserID:  "DoomGuy",
				Full:    `https://` + rndString2 + `.com`,
				Short:   rndString2,
				Deleted: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			got, err := s.GetByFull(tt.args.ctx, tt.args.fullURL)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostgres_GetByShort(t *testing.T) {
	rndString1 := keygen.GetRandkey(5)
	rndString2 := keygen.GetRandkey(5)
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	defer db.Close()
	s := &Postgres{
		db:     db,
		logger: l,
	}
	s.Store(
		_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
		domain.URL{
			UserID:  "DoomGuy",
			Full:    `https://` + rndString2 + `.com`,
			Short:   rndString2,
			Deleted: false,
		},
	)
	type fields struct {
		db     *sqlx.DB
		logger *zap.SugaredLogger
	}
	type args struct {
		ctx      _context.Context
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.URL
		wantErr bool
	}{
		{
			name: "Is not found by short - empty result",
			fields: fields{
				db:     db,
				logger: l,
			},
			args: args{
				ctx:      _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				shortURL: rndString1,
			},
			want:    domain.URL{},
			wantErr: false,
		},
		{
			name: "Is found by short",
			fields: fields{
				db:     db,
				logger: l,
			},
			args: args{
				ctx:      _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				shortURL: rndString2,
			},
			want: domain.URL{
				UserID:  "DoomGuy",
				Full:    `https://` + rndString2 + `.com`,
				Short:   rndString2,
				Deleted: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			got, err := s.GetByShort(tt.args.ctx, tt.args.shortURL)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostgres_Ping(t *testing.T) {
	ctx := _context.Background()
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	type fields struct {
		db     *sqlx.DB
		logger *zap.SugaredLogger
	}
	type args struct {
		ctx _context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Db pings",
			fields: fields{
				db:     db,
				logger: l,
			},
			args:    args{ctx: ctx},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			err := s.Ping(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestPostgres_StoreBatch(t *testing.T) {
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	type args struct {
		ctx _context.Context
		us  map[string]domain.URL
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]domain.URL
		wantErr bool
	}{
		{
			name: "Batch store items",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				us: map[string]domain.URL{
					"1": {
						UserID: "DoomGuy",
						Full:   "http://www.yandex.ru/verylongpath1",
						Short:  "short1",
					},
				},
			},
			want: map[string]domain.URL{
				"1": {
					UserID: "DoomGuy",
					Full:   "http://www.yandex.ru/verylongpath1",
					Short:  "short1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     db,
				logger: l,
			}
			got, err := s.StoreBatch(tt.args.ctx, tt.args.us)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostgres_Close(t *testing.T) {
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successful Close",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     db,
				logger: l,
			}
			err := s.Close()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPostgres_GetUserURLs(t *testing.T) {
	rndString1 := keygen.GetRandkey(5)
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	s := &Postgres{
		db:     db,
		logger: l,
	}
	s.Store(
		_context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
		domain.URL{
			UserID:  "DoomGuy",
			Full:    `https://` + rndString1 + `.com`,
			Short:   rndString1,
			Deleted: false,
		},
	)
	type args struct {
		ctx    _context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.URL
		wantErr bool
	}{
		{
			name: "User id is unkown",
			args: args{
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				userID: "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Empty user links",
			args: args{
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "Heretic"),
				userID: "Heretic",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Has user links",
			args: args{
				ctx:    _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				userID: "DoomGuy",
			},
			want: []domain.URL{
				{
					UserID:  "DoomGuy",
					Full:    `https://` + rndString1 + `.com`,
					Short:   rndString1,
					Deleted: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     db,
				logger: l,
			}
			got, err := s.GetUserURLs(tt.args.ctx, tt.args.userID)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.want != nil {
				assert.NotEmpty(t, got)
			} else {
				assert.Empty(t, got)
			}
		})
	}
}

func TestPostgres_DeleteBatch(t *testing.T) {
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	type args struct {
		ctx    _context.Context
		userID string
		pack   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "It just works, dont ask...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Postgres{
				db:     db,
				logger: l,
			}
			s.DeleteBatch(tt.args.ctx, tt.args.userID, tt.args.pack)
		})
	}
}

func TestPostgres_GetStats(t *testing.T) {
	ctx := _context.Background()
	l, _ := logger.NewLogger()
	db, _ := sqlx.Open("postgres", getDataBaseDSN())
	truncateTable(db)
	s := &Postgres{
		db:     db,
		logger: l,
	}
	tests := []struct {
		name    string
		args    func(s *Postgres)
		want    domain.Stats
		wantErr bool
	}{
		{
			name: "Stats are zeros",
			args: func(s *Postgres) {

			},
			want: domain.Stats{
				URLs:  0,
				Users: 0,
			},
			wantErr: false,
		},
		{
			name: "Get non zero stats",
			args: func(s *Postgres) {
				s.Store(
					ctx,
					domain.URL{
						UserID:  "DoomGuy",
						Full:    `https://yandex.com`,
						Short:   `short`,
						Deleted: false,
					},
				)
			},
			want: domain.Stats{
				URLs:  1,
				Users: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args(s)
			got, err := s.GetStats(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
