package filedb

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
)

type fileDBItem struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
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

func (s *FileDB) Store(ctx context.Context, u domain.URL) (domain.URL, error) {
	emptyResult := domain.URL{}
	old, err := s.GetByFull(ctx, u.Full)
	if err != nil {
		return emptyResult, nil
	}

	if old != emptyResult {
		return old, errors.ErrConflict
	}

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return emptyResult, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	item := fileDBItem{
		UUID:        uuid.NewString(),
		UserID:      u.UserID,
		ShortURL:    u.Short,
		OriginalURL: u.Full,
	}

	if err := encoder.Encode(&item); err != nil {
		return emptyResult, err
	}

	return u, nil
}

func (s *FileDB) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	return s.findInFile("OriginalURL", fullURL)
}

func (s *FileDB) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	return s.findInFile("ShortURL", shortURL)
}

func (s *FileDB) findInFile(field, needle string) (domain.URL, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return domain.URL{}, err
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

func (s *FileDB) Ping(ctx context.Context) error {
	_, error := os.Stat(s.filePath)

	return error
}

func (s *FileDB) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
	// в мапе хранится полный урл = ключ корреляции
	wantToStore := make(map[string]string, len(us))

	for k, v := range us {
		wantToStore[string(v.Full)] = k
	}

	// для начала найдем совпадения по урлу, которые были сохранены ранее
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {

		var i fileDBItem
		if err := decoder.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		k, exists := wantToStore[i.OriginalURL]
		if exists {
			// урл был сохранен ранее: удаляем из списка на сохранение и
			// восстанавливаем его старый short вместо нового
			delete(wantToStore, i.OriginalURL)
			us[k] = domain.URL{
				UserID: i.UserID,
				Full:   i.OriginalURL,
				Short:  i.ShortURL,
			}
		}

		// список на сохранение пустой, не смысла искать далее (все элементы уже есть в хранилке)
		if len(wantToStore) == 0 {
			break
		}
	}

	// все элементы на сохранение уже есть, нечего сохранять
	if len(wantToStore) == 0 {
		return us, nil
	}

	// будем сохранять только те елементы, которых нет
	file, err = os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// пишем в файл только новые
	for _, v := range wantToStore {

		item := fileDBItem{
			UUID:        uuid.NewString(),
			UserID:      us[v].UserID,
			ShortURL:    us[v].Short,
			OriginalURL: us[v].Full,
		}

		if err := encoder.Encode(&item); err != nil {
			return nil, err
		}
	}

	return us, nil
}

func (s *FileDB) GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := []domain.URL{}

	decoder := json.NewDecoder(file)
	for {
		var i fileDBItem
		if err := decoder.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if i.UserID != userID {
			continue
		}

		result = append(result, domain.URL{
			UserID: i.UserID,
			Full:   i.OriginalURL,
			Short:  i.ShortURL,
		})
	}

	return result, nil
}
