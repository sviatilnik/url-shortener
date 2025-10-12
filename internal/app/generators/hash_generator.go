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
	buf := g.pool.Get().([]byte)
	buf = buf[:0]
	buf = append(buf, hash[:]...)
	short := hex.EncodeToString(buf)
	g.pool.Put(buf)

	return short[:g.len], nil
}
