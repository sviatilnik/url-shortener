package generators

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"
)

// HashGenerator генерирует короткие коды на основе MD5-хеша входной строки.
// Использует пул объектов для оптимизации производительности.
// Генератор является потокобезопасным.
type HashGenerator struct {
	len  uint      // Длина генерируемого кода
	pool sync.Pool // Пул объектов для переиспользования буферов
}

// NewHashGenerator создает новый генератор на основе хеша.
// Параметр len определяет длину генерируемого кода.
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

// Get генерирует короткий код на основе MD5-хеша входной строки.
// Длина кода определяется при создании генератора.
// Возможные ошибки:
//   - ErrEmptyString - передана пустая строка
func (g *HashGenerator) Get(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", ErrEmptyString
	}

	hash := md5.Sum([]byte(str))
	short := hex.EncodeToString(hash[:])

	return short[:g.len], nil
}
