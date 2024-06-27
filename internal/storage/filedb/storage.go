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

type item struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type storageURL struct {
	filePath string
}

// TODO: Тут долго думал что хранить в структуре - решил хранить только путь к файлу
func NewStorageURL(fileName string) *storageURL {
	s := &storageURL{
		filePath: fileName,
	}

	return s
}

// Открываем файл на запись, пишем новый item
func (s *storageURL) Store(u domain.URL) {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Fatalf("Error opennig storage file. %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	item := item{
		UUID:        uuid.NewString(),
		ShortURL:    u.Short,
		OriginalURL: u.Full,
	}

	if err := encoder.Encode(&item); err != nil {
		logger.Log.Fatalf("Error writing struct to storage file. %v", err)
	}
}

func (s *storageURL) GetByFull(fullURL string) (domain.URL, bool) {
	return s.findInFile("OriginalURL", fullURL)
}

func (s *storageURL) GetByShort(shortURL string) (domain.URL, bool) {
	return s.findInFile("ShortURL", shortURL)
}

// TODO: Вообще это вариация на тему поскольку нужно искать совпадения в разных полях структуры
//
//		     Сделал через рефлексию, но в продовом окружении так наверное лучше не делать,
//			 искать другое решение (хотя бы банально двумя разными функциями).
//			 Еще была мысль опираться на формат данных: полный url никогда не будет похож на сокращенный
//	      	 и понимать "что в итоге нашлось short или original" на этом основании - но это попахивает наркоманией
//		     я такое не люблю :)
func (s *storageURL) findInFile(field, needle string) (domain.URL, bool) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Fatalf("Error opennig storage file. %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var i item
		if err := decoder.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			logger.Log.Fatalf("Error reading from file storage. %v", err)
		}

		if getField(&i, field) != needle {
			continue
		}

		return domain.URL{Full: i.OriginalURL, Short: i.ShortURL}, true
	}

	return domain.URL{}, false
}

func getField(i *item, field string) string {
	r := reflect.ValueOf(i)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
