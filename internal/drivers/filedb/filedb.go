// Модуль storage для хранения в файлах.
package filedb

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/keygen"
	"go.uber.org/zap"
)

type fileDBItem struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Deleted     bool   `json:"is_deleted"`
}

// Storage для хранения в файлах, включает в себя путь к файлу и логгер.
type FileDB struct {
	fileName string
	logger   *zap.SugaredLogger
}

// Конструктор storage для файла.
func NewFileDB(fileName string, logger *zap.SugaredLogger) *FileDB {
	s := &FileDB{fileName, logger}

	return s
}

// Сохранение короткой ссылки. При сохранении происходит поиск на предмет уже существующей ссылки.
// В случае если такая ссылка уже была ранее создана вернется ошибка.
func (s *FileDB) Store(ctx context.Context, u domain.URL) (domain.URL, error) {
	old, err := s.GetByFull(ctx, u.Full)
	if err != nil {
		return domain.URL{}, nil
	}

	if (old != domain.URL{}) {
		return old, errors.ErrConflict
	}

	file, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return domain.URL{}, err
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
		return domain.URL{}, err
	}

	return u, nil
}

// Поиск по полной ссылке.
func (s *FileDB) GetByFull(ctx context.Context, fullURL string) (domain.URL, error) {
	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0666)
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

		if i.OriginalURL != fullURL {
			continue
		}

		return domain.URL{
			UserID:  i.UserID,
			Full:    i.OriginalURL,
			Short:   i.ShortURL,
			Deleted: i.Deleted,
		}, nil
	}

	return domain.URL{}, nil
}

// Поиск по короткой ссылке.
func (s *FileDB) GetByShort(ctx context.Context, shortURL string) (domain.URL, error) {
	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0666)
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

		if i.ShortURL != shortURL {
			continue
		}

		return domain.URL{
			UserID:  i.UserID,
			Full:    i.OriginalURL,
			Short:   i.ShortURL,
			Deleted: i.Deleted,
		}, nil
	}

	return domain.URL{}, nil
}

// Пинг хранилки в файле.
func (s *FileDB) Ping(ctx context.Context) error {
	_, error := os.Stat(s.fileName)

	return error
}

// Пакетное сохранение коротких URL. В методе используется поиск уже существующих URL.
func (s *FileDB) StoreBatch(ctx context.Context, us map[string]domain.URL) (map[string]domain.URL, error) {
	// в мапе хранится полный урл = ключ корреляции
	wantToStore := make(map[string]string, len(us))

	for k, v := range us {
		wantToStore[string(v.Full)] = k
	}

	// для начала найдем совпадения по урлу, которые были сохранены ранее
	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {

		var i fileDBItem
		if err = decoder.Decode(&i); err == io.EOF {
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
				UserID:  i.UserID,
				Full:    i.OriginalURL,
				Short:   i.ShortURL,
				Deleted: i.Deleted,
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
	file, err = os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
			Deleted:     us[v].Deleted,
		}

		if err := encoder.Encode(&item); err != nil {
			return nil, err
		}
	}

	return us, nil
}

// Получение ссылок, созданных пользоваетелем.
func (s *FileDB) GetUserURLs(ctx context.Context, userID string) ([]domain.URL, error) {
	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]domain.URL, 0, 20)

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

// Получение рандомного ключа
func (s *FileDB) GetRandkey(n uint) string {
	return keygen.GetRandkey(n)
}
