package storages

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

type FileStorage struct {
	filePath   string
	lastUUID   int
	mut        sync.RWMutex
	cache      map[string]*models.Link
	cacheMutex sync.RWMutex
}

type storeItem struct {
	UUID        string `json:"uuid"`
	Short       string `json:"short"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
		lastUUID: 0,
		cache:    make(map[string]*models.Link),
	}
}

func (f *FileStorage) Save(ctx context.Context, link *models.Link) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

		file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		f.mut.Lock()
		defer f.mut.Unlock()

		item := &storeItem{
			OriginalURL: link.OriginalURL,
			Short:       link.ShortCode,
			UUID:        link.ID,
		}

		marshal, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		if _, err = file.Write(marshal); err != nil {
			return nil, err
		}

		file.WriteString("\n")

		// Добавляем в кэш
		f.cacheMutex.Lock()
		f.cache[link.ShortCode] = link
		f.cacheMutex.Unlock()

		f.lastUUID++
		return link, nil
	}
}

func (f *FileStorage) BatchSave(ctx context.Context, links []*models.Link) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if len(links) == 0 {
			return ErrBatchIsEmpty
		}

		for _, link := range links {
			if _, err := f.Save(ctx, link); err != nil {
				return err
			}
		}

		return nil
	}
}

func (f *FileStorage) Get(ctx context.Context, shortCode string) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Сначала проверяем кэш
		f.cacheMutex.RLock()
		if link, exists := f.cache[shortCode]; exists {
			f.cacheMutex.RUnlock()
			return link, nil
		}
		f.cacheMutex.RUnlock()

		// Если не в кэше, ищем в файле
		file, err := os.Open(f.filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			row = strings.Trim(row, "\n")

			item := &storeItem{}
			err = json.Unmarshal([]byte(row), item)
			if err != nil {
				return nil, err
			}

			if item.Short == shortCode {
				link := &models.Link{
					ID:          item.UUID,
					ShortCode:   item.Short,
					OriginalURL: item.OriginalURL,
				}

				// Добавляем в кэш
				f.cacheMutex.Lock()
				f.cache[shortCode] = link
				f.cacheMutex.Unlock()

				return link, nil
			}
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (f *FileStorage) GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error) {
	panic("implement me")
}

func (f *FileStorage) Delete(ctx context.Context, IDs []string, userID string) error {
	panic("implement me")
}
