package generators

import (
	"math/rand"
	"time"
)

type RandomGenerator struct {
	len int
}

func NewRandomGenerator(len int) *RandomGenerator {
	return &RandomGenerator{
		len: len,
	}
}

func (r RandomGenerator) Get(str string) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, r.len)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b), nil
}
