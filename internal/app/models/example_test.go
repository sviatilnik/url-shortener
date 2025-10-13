package models_test

import (
	"fmt"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

// ExampleLink демонстрирует создание и использование структуры Link.
func ExampleLink() {
	// Создаем новую ссылку
	link := models.Link{
		ID:          "unique-id-123",
		ShortCode:   "abc123",
		ShortURL:    "http://localhost:8080/abc123",
		OriginalURL: "https://example.com/very/long/url",
		UserID:      "user123",
		IsDeleted:   false,
	}

	// Выводим информацию о ссылке
	fmt.Printf("Link ID: %s\n", link.ID)
	fmt.Printf("Short Code: %s\n", link.ShortCode)
	fmt.Printf("Short URL: %s\n", link.ShortURL)
	fmt.Printf("Original URL: %s\n", link.OriginalURL)
	fmt.Printf("User ID: %s\n", link.UserID)
	fmt.Printf("Is Deleted: %t\n", link.IsDeleted)

	// Output:
	// Link ID: unique-id-123
	// Short Code: abc123
	// Short URL: http://localhost:8080/abc123
	// Original URL: https://example.com/very/long/url
	// User ID: user123
	// Is Deleted: false
}

// ExampleLink_softDelete демонстрирует мягкое удаление ссылки.
func ExampleLink_softDelete() {
	// Создаем ссылку
	link := models.Link{
		ID:          "soft-delete-test",
		ShortCode:   "del123",
		ShortURL:    "http://localhost:8080/del123",
		OriginalURL: "https://example.com/to-be-deleted",
		UserID:      "user123",
		IsDeleted:   false,
	}

	fmt.Printf("Before deletion: IsDeleted = %t\n", link.IsDeleted)

	// Выполняем мягкое удаление
	link.IsDeleted = true

	fmt.Printf("After deletion: IsDeleted = %t\n", link.IsDeleted)
	fmt.Printf("Link still exists but marked as deleted\n")

	// Output:
	// Before deletion: IsDeleted = false
	// After deletion: IsDeleted = true
	// Link still exists but marked as deleted
}

// ExampleLink_userOwnership демонстрирует принадлежность ссылки пользователю.
func ExampleLink_userOwnership() {
	// Создаем ссылки для разных пользователей
	user1Links := []models.Link{
		{
			ID:          "user1-link1",
			ShortCode:   "abc123",
			ShortURL:    "http://localhost:8080/abc123",
			OriginalURL: "https://example.com/user1-first",
			UserID:      "user1",
			IsDeleted:   false,
		},
		{
			ID:          "user1-link2",
			ShortCode:   "def456",
			ShortURL:    "http://localhost:8080/def456",
			OriginalURL: "https://example.com/user1-second",
			UserID:      "user1",
			IsDeleted:   false,
		},
	}

	user2Links := []models.Link{
		{
			ID:          "user2-link1",
			ShortCode:   "ghi789",
			ShortURL:    "http://localhost:8080/ghi789",
			OriginalURL: "https://example.com/user2-first",
			UserID:      "user2",
			IsDeleted:   false,
		},
	}

	// Выводим ссылки пользователя 1
	fmt.Printf("User1 links:\n")
	for i, link := range user1Links {
		fmt.Printf("  %d. %s -> %s\n", i+1, link.ShortCode, link.OriginalURL)
	}

	// Выводим ссылки пользователя 2
	fmt.Printf("User2 links:\n")
	for i, link := range user2Links {
		fmt.Printf("  %d. %s -> %s\n", i+1, link.ShortCode, link.OriginalURL)
	}

	// Output:
	// User1 links:
	//   1. abc123 -> https://example.com/user1-first
	//   2. def456 -> https://example.com/user1-second
	// User2 links:
	//   1. ghi789 -> https://example.com/user2-first
}

// ExampleLink_validation демонстрирует валидацию полей ссылки.
func ExampleLink_validation() {
	// Создаем валидную ссылку
	validLink := models.Link{
		ID:          "valid-id",
		ShortCode:   "valid123",
		ShortURL:    "http://localhost:8080/valid123",
		OriginalURL: "https://example.com/valid",
		UserID:      "user123",
		IsDeleted:   false,
	}

	// Проверяем валидность
	if validLink.ID != "" && validLink.ShortCode != "" && validLink.OriginalURL != "" {
		fmt.Printf("Link is valid: %s\n", validLink.ShortURL)
	}

	// Создаем невалидную ссылку
	invalidLink := models.Link{
		ID:          "", // Пустой ID
		ShortCode:   "invalid",
		ShortURL:    "http://localhost:8080/invalid",
		OriginalURL: "", // Пустой оригинальный URL
		UserID:      "user123",
		IsDeleted:   false,
	}

	// Проверяем невалидность
	if invalidLink.ID == "" || invalidLink.OriginalURL == "" {
		fmt.Printf("Link is invalid: missing required fields\n")
	}

	// Output:
	// Link is valid: http://localhost:8080/valid123
	// Link is invalid: missing required fields
}

// ExampleContextUserID демонстрирует использование константы ContextUserID.
func ExampleContextUserID() {
	// ContextUserID используется как ключ для хранения ID пользователя в контексте
	userIDKey := models.ContextUserID

	fmt.Printf("ContextUserID key type: %T\n", userIDKey)
	fmt.Printf("ContextUserID value:%v\n", userIDKey)

	// В реальном приложении это используется так:
	// ctx := context.WithValue(request.Context(), models.ContextUserID, "user123")
	// userID := ctx.Value(models.ContextUserID).(string)

	fmt.Printf("Usage: context.WithValue(ctx, models.ContextUserID, \"user123\")\n")

	// Output:
	// ContextUserID key type: models.userID
	// ContextUserID value:
	// Usage: context.WithValue(ctx, models.ContextUserID, "user123")
}

// ExampleLink_serialization демонстрирует сериализацию структуры Link.
func ExampleLink_serialization() {
	// Создаем ссылку
	link := models.Link{
		ID:          "serialization-test",
		ShortCode:   "ser123",
		ShortURL:    "http://localhost:8080/ser123",
		OriginalURL: "https://example.com/serialization",
		UserID:      "user123",
		IsDeleted:   false,
	}

	// В реальном приложении структура может быть сериализована в JSON
	// для передачи через API или сохранения в файл

	fmt.Printf("Link ready for serialization:\n")
	fmt.Printf("  ID: %s\n", link.ID)
	fmt.Printf("  ShortCode: %s\n", link.ShortCode)
	fmt.Printf("  ShortURL: %s\n", link.ShortURL)
	fmt.Printf("  OriginalURL: %s\n", link.OriginalURL)
	fmt.Printf("  UserID: %s\n", link.UserID)
	fmt.Printf("  IsDeleted: %t\n", link.IsDeleted)

	// Output:
	// Link ready for serialization:
	//   ID: serialization-test
	//   ShortCode: ser123
	//   ShortURL: http://localhost:8080/ser123
	//   OriginalURL: https://example.com/serialization
	//   UserID: user123
	//   IsDeleted: false
}
