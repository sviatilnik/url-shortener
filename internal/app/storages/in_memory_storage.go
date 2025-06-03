package storages

import (
	"context"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"strings"
)

type InMemoryStorage struct {
	store map[string]string
}

func NewInMemoryStorage() URLStorage {
	return &InMemoryStorage{
		store: make(map[string]string),
	}
}

func (i InMemoryStorage) Save(ctx context.Context, link *models.Link) (*models.Link, error) {
	if strings.TrimSpace(link.ID) == "" {
		return nil, ErrEmptyKey
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		i.store[link.ID] = link.OriginalURL
		return link, nil
	}
}

func (i InMemoryStorage) BatchSave(ctx context.Context, links []*models.Link) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for _, link := range links {
			if _, err := i.Save(ctx, link); err != nil {
				return err
			}
		}
		return nil
	}
}

func (i InMemoryStorage) Get(ctx context.Context, shortCode string) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		_, ok := i.store[shortCode]

		if !ok {
			return nil, ErrKeyNotFound
		}

		return &models.Link{
			ID:          shortCode,
			OriginalURL: i.store[shortCode],
			ShortCode:   shortCode,
		}, nil
	}
}
