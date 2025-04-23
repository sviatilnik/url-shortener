package shortener

import "errors"

var (
	ErrInvalidURL      = errors.New("invalid url")
	ErrCreateShortLink = errors.New("could not generate short link")
	ErrIDIsRequired    = errors.New("id is required")
	ErrKeyIsEmpty      = errors.New("key is empty")
	ErrKeyNotFound     = errors.New("key not found")
)
