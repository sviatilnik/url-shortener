package storages

import "errors"

type InMemoryStorage struct {
	store map[string]string
}

func NewInMemoryStorage() URLStorage {
	return &InMemoryStorage{
		store: make(map[string]string),
	}
}

func (i InMemoryStorage) Save(key string, value string) error {
	i.store[key] = value
	return nil
}

func (i InMemoryStorage) Get(key string) (string, error) {
	_, ok := i.store[key]

	if !ok {
		return "", errors.New("key not found")
	}

	return i.store[key], nil
}
