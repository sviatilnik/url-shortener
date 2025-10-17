package storages

import (
	"context"
	"strings"
	"sync"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

// InMemoryStorage представляет хранилище ссылок в памяти.
// Используется для тестирования и разработки.
// Хранилище является потокобезопасным благодаря использованию RWMutex.
type InMemoryStorage struct {
	store map[string]*models.Link // Карта для хранения коротких кодов и полных объектов Link
	mu    sync.RWMutex            // Мьютекс для обеспечения потокобезопасности
}

// NewInMemoryStorage создает новый экземпляр хранилища в памяти.
func NewInMemoryStorage() URLStorage {
	return &InMemoryStorage{
		store: make(map[string]*models.Link),
	}
}

func (i *InMemoryStorage) Save(ctx context.Context, link *models.Link) (*models.Link, error) {
	if strings.TrimSpace(link.ID) == "" {
		return nil, ErrEmptyKey
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		i.mu.Lock()
		// Создаем копию ссылки для хранения
		linkCopy := &models.Link{
			ID:          link.ID,
			ShortCode:   link.ShortCode,
			ShortURL:    link.ShortURL,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   link.IsDeleted,
		}
		i.store[link.ID] = linkCopy
		i.mu.Unlock()
		return linkCopy, nil
	}
}

func (i *InMemoryStorage) BatchSave(ctx context.Context, links []*models.Link) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if len(links) == 0 {
			return ErrBatchIsEmpty
		}
		i.mu.Lock()
		for _, link := range links {
			if strings.TrimSpace(link.ID) == "" {
				i.mu.Unlock()
				return ErrEmptyKey
			}
			// Создаем копию ссылки для хранения
			linkCopy := &models.Link{
				ID:          link.ID,
				ShortCode:   link.ShortCode,
				ShortURL:    link.ShortURL,
				OriginalURL: link.OriginalURL,
				UserID:      link.UserID,
				IsDeleted:   link.IsDeleted,
			}
			i.store[link.ID] = linkCopy
		}
		i.mu.Unlock()
		return nil
	}
}

func (i *InMemoryStorage) Get(ctx context.Context, shortCode string) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		i.mu.RLock()
		link, ok := i.store[shortCode]
		i.mu.RUnlock()

		if !ok {
			return nil, ErrKeyNotFound
		}

		// Проверяем, не удалена ли ссылка
		if link.IsDeleted {
			return nil, ErrKeyNotFound
		}

		// Возвращаем копию ссылки
		return &models.Link{
			ID:          link.ID,
			ShortCode:   link.ShortCode,
			ShortURL:    link.ShortURL,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   link.IsDeleted,
		}, nil
	}
}

func (i *InMemoryStorage) GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		i.mu.RLock()
		var userLinks []*models.Link

		for _, link := range i.store {
			// Проверяем, что ссылка принадлежит пользователю и не удалена
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
		i.mu.RUnlock()

		return userLinks, nil
	}
}

func (i *InMemoryStorage) Delete(ctx context.Context, IDs []string, userID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if len(IDs) == 0 {
			return nil
		}

		i.mu.Lock()
		for _, id := range IDs {
			if link, exists := i.store[id]; exists {
				// Проверяем, что ссылка принадлежит пользователю
				if link.UserID == userID {
					// Помечаем ссылку как удаленную (soft delete)
					link.IsDeleted = true
				}
			}
		}
		i.mu.Unlock()

		return nil
	}
}
