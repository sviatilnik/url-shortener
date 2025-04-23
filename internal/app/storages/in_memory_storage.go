package storages

import (
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

func (i InMemoryStorage) Save(key string, value string) error {
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}

	i.store[key] = value
	return nil
}

func (i InMemoryStorage) Get(key string) (string, error) {
	_, ok := i.store[key]

	if !ok {
		return "", ErrKeyNotFound
	}

	return i.store[key], nil
}
