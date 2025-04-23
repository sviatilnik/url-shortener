package generators

import (
	"math/rand"
	"strings"
	"time"
)

type RandomGenerator struct {
	len uint
}

func NewRandomGenerator(len uint) *RandomGenerator {
	return &RandomGenerator{
		len: len,
	}
}

func (r RandomGenerator) Get(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", ErrEmptyString
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, r.len)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b), nil
}
