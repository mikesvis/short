package filedb

import (
	"bufio"
	_context "context"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorageURL(t *testing.T) {
	l, _ := logger.NewLogger()
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want *FileDB
	}{
		{
			name: "New storage is of type",
			args: args{
				filePath: "dummyFile.json",
			},
			want: &FileDB{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStorage := NewFileDB(tt.args.filePath, l)
			assert.IsType(t, tt.want, newStorage)
		})
	}
}

func BenchmarkNewStorageURL(b *testing.B) {
	l, _ := logger.NewLogger()
	filepath := "dummyFile.json"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFileDB(filepath, l)
	}
}

func TestStore(t *testing.T) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	type args struct {
		item     domain.URL
		mustFail bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Store item",
			args: args{
				item: domain.URL{
					UserID: "DoomGuy",
					Full:   "http://www.yandex.ru/verylongpath",
					Short:  "short",
				},
				mustFail: false,
			},
			want: `{
				"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72",
				"user_id":"DoomGuy",
				"short_url":"short",
				"original_url":"http://www.yandex.ru/verylongpath",
				"is_deleted":false
			}`,
			wantErr: false,
		},
		{
			name: "Has conflict on store",
			args: args{
				item: domain.URL{
					UserID: "Heretic",
					Full:   "http://www.yandex.ru/verylongpath",
					Short:  "short",
				},
				mustFail: false,
			},
			want:    ``,
			wantErr: true,
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	// Creating temp file
	tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
	require.Nil(t, err)
	tmpFile.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Using temp file in storage
			fileName := tmpFile.Name()

			// Faking not existing file
			if tt.args.mustFail {
				fileName = os.TempDir() + "/!"
			}
			s := &FileDB{
				fileName: fileName,
			}

			// Storing
			result, err := s.Store(ctx, tt.args.item)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.args.item, result)

			// Reading temp file
			file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0666)
			require.Nil(t, err)
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			fileString := scanner.Text()

			assert.JSONEq(t, tt.want, fileString)

		})
	}

	// Removing temp file
	os.Remove(tmpFile.Name())
}

func BenchmarkStore(b *testing.B) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Store(ctx, domain.URL{
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath",
			Short:  "short",
		})
	}

	os.Remove(tmpFile.Name())
}

func createAndSeedTestStorage(t *testing.T, jsonString string) string {
	tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
	require.Nil(t, err)

	tmpFile.Write([]byte(jsonString))

	tmpFile.Close()
	return tmpFile.Name()
}

func TestGetByFull(t *testing.T) {
	ctx := _context.Background()

	type args struct {
		in           string
		fileContents string
	}

	tests := []struct {
		name    string
		args    args
		want    domain.URL
		wantErr bool
	}{
		{
			name: "Is found by full",
			args: args{
				in:           "http://www.yandex.ru/verylongpath",
				fileContents: `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "short",
			},
			wantErr: false,
		},
		{
			name: "Is not found by full",
			args: args{
				in:           "http://localhost/",
				fileContents: `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want:    domain.URL{},
			wantErr: false,
		},
		{
			name: "Is error",
			args: args{
				in:           "http://localhost/",
				fileContents: `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want:    domain.URL{},
			wantErr: true,
		},
		{
			name: "Is decode error",
			args: args{
				in:           "http://localhost/",
				fileContents: `-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want:    domain.URL{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createAndSeedTestStorage(t, tt.args.fileContents)
			if tt.wantErr {
				tmpFile = ""
			}
			s := &FileDB{
				fileName: tmpFile,
			}
			item, err := s.GetByFull(ctx, tt.args.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
			os.Remove(tmpFile)
		})
	}
}

func BenchmarkGetByFull(b *testing.B) {
	ctx := _context.Background()
	const JSONstring = `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Write([]byte(JSONstring))
	tmpFile.Close()
	s := &FileDB{
		fileName: tmpFile.Name(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetByFull(ctx, "http://www.yandex.ru/verylongpath")
	}

	os.Remove(tmpFile.Name())
}

func TestGetByShort(t *testing.T) {
	ctx := _context.Background()

	type args struct {
		in           string
		fileContents string
	}

	tests := []struct {
		name    string
		args    args
		want    domain.URL
		wantErr bool
	}{
		{
			name: "Is found by short",
			args: args{
				in:           "short",
				fileContents: `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "short",
			},
			wantErr: false,
		},
		{
			name: "Is not found by short",
			args: args{
				in:           "dummyShort",
				fileContents: `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`,
			},
			want:    domain.URL{},
			wantErr: false,
		},
		{
			name: "Is error",
			args: args{
				in:           "http://localhost/",
				fileContents: ``,
			},
			want:    domain.URL{},
			wantErr: true,
		},
		{
			name: "Is decode error",
			args: args{
				in:           "short",
				fileContents: `c07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://w`,
			},
			want:    domain.URL{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createAndSeedTestStorage(t, tt.args.fileContents)
			if tt.args.fileContents == "" {
				tmpFile = ""
			}
			s := &FileDB{
				fileName: tmpFile,
			}
			item, err := s.GetByShort(ctx, tt.args.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
			os.Remove(tmpFile)
		})
	}
}

func BenchmarkGetByShort(b *testing.B) {
	ctx := _context.Background()
	const JSONstring = `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Write([]byte(JSONstring))
	tmpFile.Close()
	s := &FileDB{
		fileName: tmpFile.Name(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetByShort(ctx, "short")
	}

	os.Remove(tmpFile.Name())
}

func TestStoreBatch(t *testing.T) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")

	type want struct {
		stored        map[string]domain.URL
		storageString string
	}

	tests := []struct {
		name    string
		args    map[string]domain.URL
		want    want
		wantErr bool
	}{
		{
			name: "Batch store items",
			args: map[string]domain.URL{
				"1": {
					UserID: "DoomGuy",
					Full:   "http://www.yandex.ru/verylongpath1",
					Short:  "short1",
				},
			},
			want: want{
				stored: map[string]domain.URL{
					"1": {
						UserID: "DoomGuy",
						Full:   "http://www.yandex.ru/verylongpath1",
						Short:  "short1",
					},
				},
				storageString: `{
					"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72",
					"user_id":"DoomGuy",
					"short_url": "short1",
					"original_url":"http://www.yandex.ru/verylongpath1",
					"is_deleted":false
				}`,
			},
			wantErr: false,
		},
		{
			name: "Is error",
			args: map[string]domain.URL{
				"1": {
					UserID: "DoomGuy",
					Full:   "http://www.yandex.ru/verylongpath1",
					Short:  "short1",
				},
			},
			want: want{
				stored:        map[string]domain.URL{},
				storageString: ``,
			},
			wantErr: true,
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating temp file
			tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
			require.Nil(t, err)
			tmpFile.Close()

			fileName := tmpFile.Name()

			if tt.wantErr {
				fileName = ""
			}

			// Using temp file in storage
			s := &FileDB{
				fileName: fileName,
			}

			// Storing
			stored, err := s.StoreBatch(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Reading temp file
			file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0666)
			require.Nil(t, err)
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			fileString := scanner.Text()

			assert.EqualValues(t, tt.want.stored, stored)
			assert.JSONEq(t, tt.want.storageString, fileString)

			// Removing temp file
			os.Remove(tmpFile.Name())
		})
	}
}

func BenchmarkStoreBatch(b *testing.B) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.StoreBatch(ctx, map[string]domain.URL{
			"1": {
				UserID: "DoomGuy",
				Full:   "http://www.yandex.ru/verylongpath1",
				Short:  "short1",
			},
		})
	}

	os.Remove(tmpFile.Name())
}

func TestGetUserURLs(t *testing.T) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	s.Store(ctx, domain.URL{
		UserID: "DoomGuy",
		Full:   "http://iddqd.com",
		Short:  "idkfa",
	})

	tests := []struct {
		name string
		args string
		want []domain.URL
	}{
		{
			name: "Get current user URLs",
			args: "DoomGuy",
			want: []domain.URL{
				{
					UserID:  "DoomGuy",
					Full:    "http://iddqd.com",
					Short:   "idkfa",
					Deleted: false,
				},
			},
		},
		{
			name: "Get empty list user URLs",
			args: "Heretic",
			want: []domain.URL{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.GetUserURLs(ctx, tt.args)
			require.NoError(t, err)

			assert.EqualValues(t, tt.want, result)
		})
	}
	os.Remove(tmpFile.Name())
}

func BenchmarkGetUserURLs(b *testing.B) {
	ctx := _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy")
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	s.Store(ctx, domain.URL{
		UserID: "DoomGuy",
		Full:   "http://iddqd.com",
		Short:  "idkfa",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetUserURLs(ctx, "DoomGuy")
	}

	os.Remove(tmpFile.Name())
}

func TestFileDB_Ping(t *testing.T) {
	ctx := _context.Background()
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}

	tests := []struct {
		name    string
		args    _context.Context
		wantErr bool
	}{
		{
			name:    "Ping storage",
			args:    ctx,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Ping(tt.args)
			require.NoError(t, err)
		})
	}

	os.Remove(tmpFile.Name())
}

func TestFileDB_GetRandkey(t *testing.T) {
	l, _ := logger.NewLogger()
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
			s := &FileDB{
				fileName: "/tmp/dummyFile.db",
				logger:   l,
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

func TestGetStats(t *testing.T) {
	ctx := _context.Background()
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Close()

	s := &FileDB{
		fileName: tmpFile.Name(),
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))

	tests := []struct {
		name    string
		args    func(s *FileDB)
		want    domain.Stats
		wantErr bool
	}{
		{
			name: "Get zeros",
			args: func(s *FileDB) {},
			want: domain.Stats{
				URLs:  0,
				Users: 0,
			},
			wantErr: false,
		},
		{
			name: "Get non zero stats",
			args: func(s *FileDB) {
				s.Store(ctx, domain.URL{
					UserID: "DoomGuy1",
					Full:   "http://iddqd1.com",
					Short:  "idkfa1",
				})
				s.Store(ctx, domain.URL{
					UserID: "DoomGuy2",
					Full:   "http://iddqd2.com",
					Short:  "idkfa2",
				})
			},
			want: domain.Stats{
				URLs:  2,
				Users: 2,
			},
			wantErr: false,
		},
		{
			name: "Get error on corrupted file",
			args: func(s *FileDB) {
				f, _ := os.OpenFile(s.fileName, os.O_RDWR, 0644)
				f.Truncate(50)
			},
			want:    domain.Stats{},
			wantErr: true,
		},
		{
			name: "Filename is wrong",
			args: func(s *FileDB) {
				s.fileName = ""
			},
			want:    domain.Stats{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args(s)
			result, err := s.GetStats(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, result)
		})
	}
	os.Remove(tmpFile.Name())
}
