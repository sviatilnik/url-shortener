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
	UserID      string `json:"user_id"`
	IsDeleted   bool   `json:"is_deleted"`
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
			UserID:      link.UserID,
			IsDeleted:   link.IsDeleted,
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
			// Проверяем, не удалена ли ссылка
			if link.IsDeleted {
				return nil, ErrKeyNotFound
			}
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
				// Проверяем, не удалена ли ссылка
				if item.IsDeleted {
					return nil, ErrKeyNotFound
				}

				link := &models.Link{
					ID:          item.UUID,
					ShortCode:   item.Short,
					OriginalURL: item.OriginalURL,
					UserID:      item.UserID,
					IsDeleted:   item.IsDeleted,
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var userLinks []*models.Link

		// Сначала проверяем кэш
		f.cacheMutex.RLock()
		for _, link := range f.cache {
			if link.UserID == userID && !link.IsDeleted {
				// Создаем копию ссылки
				linkCopy := &models.Link{
					ID:          link.ID,
					ShortCode:   link.ShortCode,
					ShortURL:    link.ShortURL,
					OriginalURL: link.OriginalURL,
					UserID:      link.UserID,
					IsDeleted:   link.IsDeleted,
				}
				userLinks = append(userLinks, linkCopy)
			}
		}
		f.cacheMutex.RUnlock()

		// Если файл не существует, возвращаем то, что есть в кэше
		file, err := os.Open(f.filePath)
		if err != nil {
			return userLinks, nil
		}
		defer file.Close()

		// Сканируем файл для поиска ссылок, которых нет в кэше
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			row = strings.Trim(row, "\n")

			item := &storeItem{}
			err = json.Unmarshal([]byte(row), item)
			if err != nil {
				continue // Пропускаем некорректные записи
			}

			// Проверяем, что ссылка принадлежит пользователю и не удалена
			if item.UserID == userID && !item.IsDeleted {
				// Проверяем, есть ли уже в кэше
				f.cacheMutex.RLock()
				_, exists := f.cache[item.Short]
				f.cacheMutex.RUnlock()

				if !exists {
					link := &models.Link{
						ID:          item.UUID,
						ShortCode:   item.Short,
						OriginalURL: item.OriginalURL,
						UserID:      item.UserID,
						IsDeleted:   item.IsDeleted,
					}
					userLinks = append(userLinks, link)
				}
			}
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		return userLinks, nil
	}
}

func (f *FileStorage) Delete(ctx context.Context, IDs []string, userID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if len(IDs) == 0 {
			return nil
		}

		// Создаем множество ID для быстрого поиска
		idsToDelete := make(map[string]bool)
		for _, id := range IDs {
			idsToDelete[id] = true
		}

		// Обновляем кэш
		f.cacheMutex.Lock()
		for _, link := range f.cache {
			if idsToDelete[link.ID] && link.UserID == userID {
				link.IsDeleted = true
			}
		}
		f.cacheMutex.Unlock()

		// Читаем все записи из файла
		file, err := os.Open(f.filePath)
		if err != nil {
			// Если файл не существует, возвращаем nil
			return nil
		}
		defer file.Close()

		var items []*storeItem
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			row = strings.Trim(row, "\n")

			item := &storeItem{}
			err = json.Unmarshal([]byte(row), item)
			if err != nil {
				continue // Пропускаем некорректные записи
			}

			// Помечаем как удаленные ссылки, принадлежащие пользователю
			if idsToDelete[item.UUID] && item.UserID == userID {
				item.IsDeleted = true
			}

			items = append(items, item)
		}

		err = scanner.Err()
		if err != nil {
			return err
		}

		// Перезаписываем файл с обновленными данными
		tempFile, err := os.CreateTemp("", "url_shortener_*.tmp")
		if err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())

		for _, item := range items {
			marshal, err := json.Marshal(item)
			if err != nil {
				tempFile.Close()
				return err
			}

			if _, err = tempFile.Write(marshal); err != nil {
				tempFile.Close()
				return err
			}

			if _, err = tempFile.WriteString("\n"); err != nil {
				tempFile.Close()
				return err
			}
		}

		tempFile.Close()

		// Атомарно заменяем оригинальный файл
		return os.Rename(tempFile.Name(), f.filePath)
	}
}
