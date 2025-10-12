package generators

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"
)

type HashGenerator struct {
	len  uint
	pool sync.Pool
}

func NewHashGenerator(len uint) *HashGenerator {
	return &HashGenerator{
		len: len,
		pool: sync.Pool{
			New: func() any {
				return make([]byte, 0, 32)
			},
		},
	}
}

func (g *HashGenerator) Get(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", ErrEmptyString
	}

	hash := md5.Sum([]byte(str))
	short := hex.EncodeToString(hash[:])

	return short[:g.len], nil
}
