package config

import (
	"errors"
	"strings"
)

type StaticConfig struct {
	store map[string]interface{}
}

func NewStaticConfig(store map[string]interface{}) StaticConfig {
	if store == nil {
		store = make(map[string]interface{})
	}

	return StaticConfig{
		store: store,
	}
}

func (s StaticConfig) Get(key string, default_value interface{}) interface{} {
	if v, ok := s.store[key]; ok {
		return v
	}
	return default_value
}

func (s StaticConfig) Set(key string, value interface{}) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key is empty")
	}
	s.store[key] = value
	return nil
}
