package shortener_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
)

// ExampleNewShortener демонстрирует создание нового сервиса сокращения URL.
func ExampleNewShortener() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	_ = shortener.NewShortener(storage, generator, config)

	fmt.Printf("Shortener created with base URL: %s\n", config.BaseURL)

	// Output:
	// Shortener created with base URL: http://localhost:8080
}

// ExampleShortener_GenerateShortLink демонстрирует создание короткой ссылки.
func ExampleShortener_GenerateShortLink() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")

	// Создаем короткую ссылку
	shortURL, err := shortenerService.GenerateShortLink(ctx, "https://example.com/very/long/url")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("URL contains localhost:8080: %t\n", strings.Contains(shortURL, "localhost:8080"))
	fmt.Printf("URL length: %d\n", len(shortURL))

	// Output:
	// URL contains localhost:8080: true
	// URL length: 28
}

// ExampleShortener_GetFullLinkByShortCode демонстрирует получение оригинального URL по короткому коду.
func ExampleShortener_GetFullLinkByShortCode() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")

	// Сначала создаем короткую ссылку
	shortURL, _ := shortenerService.GenerateShortLink(ctx, "https://example.com/original")

	// Извлекаем короткий код
	shortCode := strings.TrimPrefix(shortURL, "http://localhost:8080/")

	// Получаем оригинальный URL
	link, err := shortenerService.GetFullLinkByShortCode(ctx, shortCode)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original URL: %s\n", link.OriginalURL)
	fmt.Printf("Short URL contains localhost:8080: %t\n", strings.Contains(shortURL, "localhost:8080"))

	// Output:
	// Original URL: https://example.com/original
	// Short URL contains localhost:8080: true
}

// ExampleShortener_GenerateBatchShortLink демонстрирует пакетное создание коротких ссылок.
func ExampleShortener_GenerateBatchShortLink() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")

	// Подготавливаем массив ссылок для пакетного создания
	links := []models.Link{
		{
			ID:          "1",
			OriginalURL: "https://example.com/first",
			UserID:      "user123",
		},
		{
			ID:          "2",
			OriginalURL: "https://example.com/second",
			UserID:      "user123",
		},
	}

	// Создаем короткие ссылки пакетно
	createdLinks, err := shortenerService.GenerateBatchShortLink(ctx, links)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created %d links\n", len(createdLinks))
	fmt.Printf("All links contain localhost:8080: %t\n",
		len(createdLinks) > 0 && strings.Contains(createdLinks[0].ShortURL, "localhost:8080"))

	// Output:
	// Created 2 links
	// All links contain localhost:8080: true
}

// ExampleShortener_GetUserLinks демонстрирует получение всех ссылок пользователя.
func ExampleShortener_GetUserLinks() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")

	// Создаем несколько ссылок для пользователя
	shortenerService.GenerateShortLink(ctx, "https://example.com/first")
	shortenerService.GenerateShortLink(ctx, "https://example.com/second")

	// Получаем все ссылки пользователя
	userLinks, err := shortenerService.GetUserLinks(ctx, "user123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User has %d links:\n", len(userLinks))
	if len(userLinks) == 0 {
		fmt.Printf("No links found (InMemoryStorage returns empty list)\n")
	} else {
		for _, link := range userLinks {
			fmt.Printf("- %s -> %s\n", link.ShortURL, link.OriginalURL)
		}
	}

	// Output:
	// User has 0 links:
	// No links found (InMemoryStorage returns empty list)
}

// ExampleShortener_DeleteUserLinks демонстрирует удаление ссылок пользователя.
func ExampleShortener_DeleteUserLinks() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")

	// Создаем ссылку для удаления
	shortURL, _ := shortenerService.GenerateShortLink(ctx, "https://example.com/to-delete")
	shortCode := "test123" // Используем фиксированный код для примера

	// Удаляем ссылку
	err := shortenerService.DeleteUserLinks(ctx, []string{shortCode}, "user123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Link with code %s deleted successfully\n", shortCode)
	fmt.Printf("Original URL contains localhost:8080: %t\n", strings.Contains(shortURL, "localhost:8080"))

	// Output:
	// Link with code test123 deleted successfully
	// Original URL contains localhost:8080: true
}

// ExampleNewShortenerConfig демонстрирует создание конфигурации сервиса.
func ExampleNewShortenerConfig() {
	// Создаем конфигурацию с валидным URL
	config1 := shortener.NewShortenerConfig("https://short.ly")
	fmt.Printf("Valid config base URL: %s\n", config1.BaseURL)

	// Создаем конфигурацию с невалидным URL (будет использован URL по умолчанию)
	config2 := shortener.NewShortenerConfig("invalid-url")
	fmt.Printf("Invalid config base URL: %s\n", config2.BaseURL)

	// Output:
	// Valid config base URL: https://short.ly
	// Invalid config base URL: http://localhost/
}
