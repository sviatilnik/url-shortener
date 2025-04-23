package config

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"strings"
)

type DefaultConfig struct {
	store map[string]interface{}
}

func NewDefaultConfig(store map[string]interface{}) Config {
	if store == nil {
		store = make(map[string]interface{})
	}

	conf := &DefaultConfig{store: store}
	conf.SetDefaultValues()

	return conf
}

func (s *DefaultConfig) Get(key string) interface{} {
	if v, ok := s.store[key]; ok {
		return v
	}
	return nil
}

func (s *DefaultConfig) Set(key string, value interface{}) error {
	if strings.TrimSpace(key) == "" {
		return shortener.ErrKeyIsEmpty
	}
	s.store[key] = value
	return nil
}

func (s *DefaultConfig) SetDefaultValues() {
	s.store["host"] = "localhost:8080"
	s.store["shortURLHost"] = "http://localhost:8080"
}
