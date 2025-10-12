package generators

import (
	"math/rand"
	"strings"
	"sync"
)

type RandomGenerator struct {
	len uint
	rnd *rand.Rand
	mu  sync.Mutex
}

func NewRandomGenerator(len uint) *RandomGenerator {
	return &RandomGenerator{
		len: len,
		rnd: rand.New(rand.NewSource(rand.Int63())),
	}
}

func (r *RandomGenerator) Get(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", ErrEmptyString
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, r.len)
	for i := range b {
		b[i] = chars[r.rnd.Intn(len(chars))]
	}

	return string(b), nil
}
