package generators

import (
	"math/rand"
	"strings"
	"sync"
)

// RandomGenerator генерирует случайные короткие коды заданной длины.
// Использует алфавит из букв (A-Z, a-z) и цифр (0-9).
// Генератор является потокобезопасным.
type RandomGenerator struct {
	len uint       // Длина генерируемого кода
	rnd *rand.Rand // Генератор случайных чисел
	mu  sync.Mutex // Мьютекс для обеспечения потокобезопасности
}

// NewRandomGenerator создает новый генератор случайных кодов.
// Параметр len определяет длину генерируемого кода.
func NewRandomGenerator(len uint) *RandomGenerator {
	return &RandomGenerator{
		len: len,
		rnd: rand.New(rand.NewSource(rand.Int63())),
	}
}

// Get генерирует случайный короткий код для переданной строки.
// Длина кода определяется при создании генератора.
// Возможные ошибки:
//   - ErrEmptyString - передана пустая строка
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
