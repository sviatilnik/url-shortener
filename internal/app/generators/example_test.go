package generators_test

import (
	"fmt"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
)

// isAlphanumeric проверяет, содержит ли строка только буквы и цифры
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// isHex проверяет, содержит ли строка только шестнадцатеричные символы
func isHex(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

// ExampleNewRandomGenerator демонстрирует создание генератора случайных кодов.
func ExampleNewRandomGenerator() {
	// Создаем генератор с длиной кода 6 символов
	generator := generators.NewRandomGenerator(6)

	// Генерируем код для URL
	code, err := generator.Get("https://example.com/very/long/url")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Code length: %d\n", len(code))
	fmt.Printf("Code contains only alphanumeric chars: %t\n", isAlphanumeric(code))

	// Output:
	// Code length: 6
	// Code contains only alphanumeric chars: true
}

// ExampleRandomGenerator_Get демонстрирует генерацию случайных кодов.
func ExampleRandomGenerator_Get() {
	// Создаем генератор
	generator := generators.NewRandomGenerator(8)

	// Генерируем коды для разных URL
	urls := []string{
		"https://example.com/first",
		"https://example.com/second",
		"https://example.com/third",
	}

	for i, url := range urls {
		code, err := generator.Get(url)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", url, err)
			continue
		}
		fmt.Printf("URL %d: code length %d, alphanumeric: %t\n", i+1, len(code), isAlphanumeric(code))
	}

	// Output:
	// URL 1: code length 8, alphanumeric: true
	// URL 2: code length 8, alphanumeric: true
	// URL 3: code length 8, alphanumeric: true
}

// ExampleNewHashGenerator демонстрирует создание генератора на основе хеша.
func ExampleNewHashGenerator() {
	// Создаем генератор с длиной кода 8 символов
	generator := generators.NewHashGenerator(8)

	// Генерируем код для URL
	code, err := generator.Get("https://example.com/very/long/url")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Code length: %d\n", len(code))
	fmt.Printf("Code contains only hex chars: %t\n", isHex(code))

	// Output:
	// Code length: 8
	// Code contains only hex chars: true
}

// ExampleHashGenerator_Get демонстрирует генерацию кодов на основе хеша.
func ExampleHashGenerator_Get() {
	// Создаем генератор
	generator := generators.NewHashGenerator(6)

	// Генерируем коды для разных URL
	urls := []string{
		"https://example.com/first",
		"https://example.com/second",
		"https://example.com/third",
	}

	for i, url := range urls {
		code, err := generator.Get(url)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", url, err)
			continue
		}
		fmt.Printf("URL %d: code %s length %d, hex: %t\n", i+1, code, len(code), isHex(code))
	}

	// Output:
	// URL 1: code length 6, hex: true
	// URL 2: code length 6, hex: true
	// URL 3: code length 6, hex: true
}
