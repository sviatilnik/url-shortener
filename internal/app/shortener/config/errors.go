package config

import "errors"

var (
	ErrInvalidURL  = errors.New("invalid url")
	ErrKeyIsEmpty  = errors.New("key is empty")
	ErrKeyNotFound = errors.New("key not found")
)
