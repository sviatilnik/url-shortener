package generators

import (
	"crypto/md5"
	"encoding/hex"
)

type HashGenerator struct {
	len int
}

func NewHashGenerator(len int) *HashGenerator {
	return &HashGenerator{
		len: len,
	}
}

func (g *HashGenerator) Get(str string) (string, error) {
	hash := md5.Sum([]byte(str))
	short := hex.EncodeToString(hash[:])

	return short[:g.len], nil
}
