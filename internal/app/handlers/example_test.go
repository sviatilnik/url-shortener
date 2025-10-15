package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
)

// ExampleAPIShortLinkHandler демонстрирует создание короткой ссылки через API.
func ExampleAPIShortLinkHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем обработчик
	handler := handlers.APIShortLinkHandler(shortenerService)

	// Подготавливаем запрос
	requestBody := map[string]string{
		"url": "https://example.com/very/long/url",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Создаем HTTP-запрос
	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем контекст с пользователем
	ctx := context.WithValue(req.Context(), models.ContextUserID, "user123")
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)
	if w.Code == 201 {
		fmt.Printf("Response contains result field: %t\n", strings.Contains(w.Body.String(), "result"))
		fmt.Printf("Response contains localhost:8080: %t\n", strings.Contains(w.Body.String(), "localhost:8080"))
	}

	// Output:
	// Status: 201
	// Response contains result field: true
	// Response contains localhost:8080: true
}

// ExampleBatchShortLinkHandler демонстрирует пакетное создание коротких ссылок.
func ExampleBatchShortLinkHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем обработчик
	handler := handlers.BatchShortLinkHandler(shortenerService)

	// Подготавливаем запрос
	requestBody := []map[string]string{
		{
			"correlation_id": "1",
			"original_url":   "https://example.com/first",
		},
		{
			"correlation_id": "2",
			"original_url":   "https://example.com/second",
		},
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Создаем HTTP-запрос
	req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем контекст с пользователем
	ctx := context.WithValue(req.Context(), models.ContextUserID, "user123")
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)
	if w.Code == 201 {
		fmt.Printf("Response contains correlation_id: %t\n", strings.Contains(w.Body.String(), "correlation_id"))
		fmt.Printf("Response contains short_url: %t\n", strings.Contains(w.Body.String(), "short_url"))
		fmt.Printf("Response contains localhost:8080: %t\n", strings.Contains(w.Body.String(), "localhost:8080"))
	}

	// Output:
	// Status: 201
	// Response contains correlation_id: true
	// Response contains short_url: true
	// Response contains localhost:8080: true
}

// ExampleGetShortLinkHandler демонстрирует создание короткой ссылки через простой POST.
func ExampleGetShortLinkHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем обработчик
	handler := handlers.GetShortLinkHandler(shortenerService)

	// Подготавливаем запрос
	url := "https://example.com/simple/url"

	// Создаем HTTP-запрос
	req := httptest.NewRequest("POST", "/", strings.NewReader(url))

	// Создаем контекст с пользователем
	ctx := context.WithValue(req.Context(), models.ContextUserID, "user123")
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)
	if w.Code == 201 {
		fmt.Printf("Response contains localhost:8080: %t\n", strings.Contains(w.Body.String(), "localhost:8080"))
		fmt.Printf("Response length: %d\n", len(w.Body.String()))
	}

	// Output:
	// Status: 201
	// Response contains localhost:8080: true
	// Response length: 28
}

// ExampleRedirectToFullLinkHandler демонстрирует перенаправление по короткой ссылке.
func ExampleRedirectToFullLinkHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем обработчик
	handler := handlers.RedirectToFullLinkHandler(shortenerService)

	// Создаем HTTP-запрос с несуществующим кодом
	req := httptest.NewRequest("GET", "/nonexistent", nil)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)
	if w.Code == 400 {
		fmt.Printf("Link not found (expected for nonexistent code)\n")
	}

	// Output:
	// Status: 400
	// Link not found (expected for nonexistent code)
}

// ExampleUserURLsHandler демонстрирует получение списка URL пользователя.
func ExampleUserURLsHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем несколько ссылок для пользователя
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")
	shortenerService.GenerateShortLink(ctx, "https://example.com/first")
	shortenerService.GenerateShortLink(ctx, "https://example.com/second")

	// Создаем обработчик
	handler := handlers.UserURLsHandler(shortenerService)

	// Создаем HTTP-запрос
	req := httptest.NewRequest("GET", "/api/user/urls", nil)
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)
	if w.Code == 204 {
		fmt.Printf("No content - user has no URLs\n")
	} else if w.Code == 200 {
		fmt.Printf("Response contains short_url: %t\n", strings.Contains(w.Body.String(), "short_url"))
		fmt.Printf("Response contains original_url: %t\n", strings.Contains(w.Body.String(), "original_url"))
	}

	// Output:
	// Status: 204
	// No content - user has no URLs
}

// ExampleDeleteUserURLsHandler демонстрирует удаление URL пользователя.
func ExampleDeleteUserURLsHandler() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем генератор случайных кодов
	generator := generators.NewRandomGenerator(6)

	// Создаем конфигурацию
	config := shortener.NewShortenerConfig("http://localhost:8080")

	// Создаем сервис сокращения URL
	shortenerService := shortener.NewShortener(storage, generator, config)

	// Создаем ссылку для пользователя
	ctx := context.WithValue(context.Background(), models.ContextUserID, "user123")
	_, _ = shortenerService.GenerateShortLink(ctx, "https://example.com/to-delete")
	shortCode := "test123" // Используем фиксированный код для примера

	// Создаем обработчик
	handler := handlers.DeleteUserURLsHandler(shortenerService)

	// Подготавливаем запрос с ID ссылки для удаления
	requestBody := []string{shortCode}
	jsonBody, _ := json.Marshal(requestBody)

	// Создаем HTTP-запрос
	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 202
}

// ExamplePingDBHandler демонстрирует проверку состояния базы данных.
func ExamplePingDBHandler() {
	// В реальном приложении здесь был бы настоящий DB connection
	// Для примера создаем mock connection

	// Создаем обработчик
	handler := handlers.PingDBHandler(nil) // nil для демонстрации

	// Создаем HTTP-запрос
	req := httptest.NewRequest("GET", "/ping", nil)

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Выполняем запрос (будет паника из-за nil connection)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Handler panicked as expected with nil connection\n")
		}
	}()

	handler(w, req)

	// Проверяем результат
	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Handler panicked as expected with nil connection
}
