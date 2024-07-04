package filedb

import (
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/logger"
)

type fileDbItem struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileDb struct {
	filePath string
}

func NewFileDb(fileName string) *FileDb {
	s := &FileDb{
		filePath: fileName,
	}

	return s
}

func (s *FileDb) Store(u domain.URL) {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Fatalf("Error opennig storage file. %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	item := fileDbItem{
		UUID:        uuid.NewString(),
		ShortURL:    u.Short,
		OriginalURL: u.Full,
	}

	if err := encoder.Encode(&item); err != nil {
		logger.Log.Fatalf("Error writing struct to storage file. %v", err)
	}
}

func (s *FileDb) GetByFull(fullURL string) (domain.URL, error) {
	return s.findInFile("OriginalURL", fullURL)
}

func (s *FileDb) GetByShort(shortURL string) (domain.URL, error) {
	return s.findInFile("ShortURL", shortURL)
}

func (s *FileDb) findInFile(field, needle string) (domain.URL, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Fatalf("Error opennig storage file. %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var i fileDbItem
		if err := decoder.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			return domain.URL{}, err
		}

		if getField(&i, field) != needle {
			continue
		}

		return domain.URL{Full: i.OriginalURL, Short: i.ShortURL}, nil
	}

	return domain.URL{}, nil
}

func getField(i *fileDbItem, field string) string {
	r := reflect.ValueOf(i)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
