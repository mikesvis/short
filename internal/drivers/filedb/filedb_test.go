package filedb

import (
	"bufio"
	"context"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorageURL(t *testing.T) {
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
			newStorage := NewFileDB(tt.args.filePath)
			assert.IsType(t, tt.want, newStorage)
		})
	}
}

func Test_storageURL_Store(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		args domain.URL
		want string
	}{
		{
			name: "Store item",
			args: domain.URL{
				Full:  "http://www.yandex.ru/verylongpath",
				Short: "short",
			},
			want: `{
				"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72",
				"short_url": "short",
				"original_url":"http://www.yandex.ru/verylongpath"
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
				filePath: tmpFile.Name(),
			}

			// Storing
			s.Store(ctx, tt.args)

			// Reading temp file
			file, err := os.OpenFile(s.filePath, os.O_RDONLY, 0666)
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

func createAndSeedTestStorage(t *testing.T) string {
	const JSONstring = `{"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72","short_url": "short","original_url":"http://www.yandex.ru/verylongpath"}`
	tmpFile, err := os.CreateTemp(os.TempDir(), "dbtest*.json")
	require.Nil(t, err)

	tmpFile.Write([]byte(JSONstring))

	tmpFile.Close()
	return tmpFile.Name()
}

func Test_storageURL_GetByFull(t *testing.T) {
	ctx := context.Background()

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
				filePath: tmpFile,
			}
			item, _ := s.GetByFull(ctx, tt.args)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
			os.Remove(tmpFile)
		})
	}
}

func Test_storageURL_GetByShort(t *testing.T) {
	ctx := context.Background()

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
				filePath: tmpFile,
			}
			item, _ := s.GetByShort(ctx, tt.args)
			assert.IsType(t, tt.want, item)
			assert.EqualValues(t, tt.want, item)
			os.Remove(tmpFile)
		})
	}
}

func TestFileDB_StoreBatch(t *testing.T) {
	ctx := context.Background()

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
					Full:  "http://www.yandex.ru/verylongpath1",
					Short: "short1",
				},
			},
			want: want{
				stored: map[string]domain.URL{
					"1": {
						Full:  "http://www.yandex.ru/verylongpath1",
						Short: "short1",
					},
				},
				storageString: `{
					"uuid":"52fdfc07-2182-454f-963f-5f0f9a621d72",
					"short_url": "short1",
					"original_url":"http://www.yandex.ru/verylongpath1"
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
				filePath: tmpFile.Name(),
			}

			// Storing
			stored, err := s.StoreBatch(ctx, tt.args)
			require.NoError(t, err)

			// Reading temp file
			file, err := os.OpenFile(s.filePath, os.O_RDONLY, 0666)
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
