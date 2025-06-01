package storages

import (
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

func (i InMemoryStorage) Save(link *models.Link) error {
	if strings.TrimSpace(link.ID) == "" {
		return ErrEmptyKey
	}

	i.store[link.ID] = link.OriginalURL
	return nil
}

func (i InMemoryStorage) BatchSave(links []*models.Link) error {
	for _, link := range links {
		if err := i.Save(link); err != nil {
			return err
		}
	}

	return nil
}

func (i InMemoryStorage) Get(shortCode string) (*models.Link, error) {
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
