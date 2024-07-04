package filedb

import (
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
)

type fileDBItem struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileDB struct {
	filePath string
}

func NewFileDB(fileName string) *FileDB {
	s := &FileDB{
		filePath: fileName,
	}

	return s
}

func (s *FileDB) Store(u domain.URL) error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	item := fileDBItem{
		UUID:        uuid.NewString(),
		ShortURL:    u.Short,
		OriginalURL: u.Full,
	}

	if err := encoder.Encode(&item); err != nil {
		return err
	}

	return nil
}

func (s *FileDB) GetByFull(fullURL string) (domain.URL, error) {
	return s.findInFile("OriginalURL", fullURL)
}

func (s *FileDB) GetByShort(shortURL string) (domain.URL, error) {
	return s.findInFile("ShortURL", shortURL)
}

func (s *FileDB) findInFile(field, needle string) (domain.URL, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return domain.URL{}, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var i fileDBItem
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

func getField(i *fileDBItem, field string) string {
	r := reflect.ValueOf(i)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
