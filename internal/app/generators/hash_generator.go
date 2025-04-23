package generators

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type HashGenerator struct {
	len uint
}

func NewHashGenerator(len uint) *HashGenerator {
	return &HashGenerator{
		len: len,
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
