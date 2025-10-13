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
	store map[string]string // Карта для хранения коротких кодов и оригинальных URL
	mu    sync.RWMutex      // Мьютекс для обеспечения потокобезопасности
}

// NewInMemoryStorage создает новый экземпляр хранилища в памяти.
func NewInMemoryStorage() URLStorage {
	return &InMemoryStorage{
		store: make(map[string]string),
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
		i.store[link.ID] = link.OriginalURL
		i.mu.Unlock()
		return link, nil
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
			i.store[link.ID] = link.OriginalURL
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
		originalURL, ok := i.store[shortCode]
		i.mu.RUnlock()

		if !ok {
			return nil, ErrKeyNotFound
		}

		return &models.Link{
			ID:          shortCode,
			OriginalURL: originalURL,
			ShortCode:   shortCode,
		}, nil
	}
}

func (i *InMemoryStorage) GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Для простоты возвращаем пустой список
		// В реальной реализации нужно хранить информацию о пользователях
		return []*models.Link{}, nil
	}
}

func (i *InMemoryStorage) Delete(ctx context.Context, IDs []string, userID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Для простоты просто возвращаем nil
		// В реальной реализации нужно помечать ссылки как удаленные
		return nil
	}
}
