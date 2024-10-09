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

	tests := []struct {
		name string
		args domain.URL
		want string
	}{
		{
			name: "Store item",
			args: domain.URL{
				UserID: "DoomGuy",
				Full:   "http://www.yandex.ru/verylongpath",
				Short:  "short",
			},
			want: `{
				"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72",
				"user_id":"DoomGuy",
				"short_url":"short",
				"original_url":"http://www.yandex.ru/verylongpath",
				"is_deleted":false
			}`,
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating temp file
			tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
			require.Nil(t, err)
			tmpFile.Close()

			// Using temp file in storage
			s := &FileDB{
				fileName: tmpFile.Name(),
			}

			// Storing
			result, err := s.Store(ctx, tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.args, result)

			// Reading temp file
			file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0666)
			require.Nil(t, err)
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			fileString := scanner.Text()

			assert.JSONEq(t, tt.want, fileString)

			// Removing temp file
			os.Remove(tmpFile.Name())
		})
	}
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

func createAndSeedTestStorage(t *testing.T) string {
	const JSONstring = `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`
	tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
	require.Nil(t, err)

	tmpFile.Write([]byte(JSONstring))

	tmpFile.Close()
	return tmpFile.Name()
}

func TestGetByFull(t *testing.T) {
	ctx := _context.Background()

	tests := []struct {
		name string
		args string
		want domain.URL
	}{
		{
			name: "Is found by full",
			args: "http://www.yandex.ru/verylongpath",
			want: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "short",
			},
		}, {
			name: "Is not found by full",
			args: "http://localhost/",
			want: domain.URL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createAndSeedTestStorage(t)
			s := &FileDB{
				fileName: tmpFile,
			}
			item, _ := s.GetByFull(ctx, tt.args)
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

	tests := []struct {
		name string
		args string
		want domain.URL
	}{
		{
			name: "Is found by short",
			args: "short",
			want: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "short",
			},
		}, {
			name: "Is not found by short",
			args: "dummyShort",
			want: domain.URL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createAndSeedTestStorage(t)
			s := &FileDB{
				fileName: tmpFile,
			}
			item, _ := s.GetByShort(ctx, tt.args)
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
		name string
		args map[string]domain.URL
		want want
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
		},
	}
	uuid.SetRand(rand.New(rand.NewSource(1)))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating temp file
			tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
			require.Nil(t, err)
			tmpFile.Close()

			// Using temp file in storage
			s := &FileDB{
				fileName: tmpFile.Name(),
			}

			// Storing
			stored, err := s.StoreBatch(ctx, tt.args)
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
