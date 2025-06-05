package storages

import "errors"

var (
	ErrKeyNotFound              = errors.New("key not found")
	ErrEmptyKey                 = errors.New("empty key")
	ErrOriginalURLAlreadyExists = errors.New("original url already exists")
	ErrBatchIsEmpty             = errors.New("batch is empty")
)
