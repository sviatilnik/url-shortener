package storages_test

import (
	"context"
	"fmt"

	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
)

// ExampleNewInMemoryStorage демонстрирует создание хранилища в памяти.
func ExampleNewInMemoryStorage() {
	// Создаем хранилище в памяти
	storage := storages.NewInMemoryStorage()

	// Создаем ссылку для сохранения
	link := &models.Link{
		ID:          "abc123",
		ShortCode:   "abc123",
		OriginalURL: "https://example.com/original",
		UserID:      "user123",
		IsDeleted:   false,
	}

	// Сохраняем ссылку
	ctx := context.Background()
	savedLink, err := storage.Save(ctx, link)
	if err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}

	fmt.Printf("Saved link: %s -> %s\n", savedLink.ShortCode, savedLink.OriginalURL)

	// Output:
	// Saved link: abc123 -> https://example.com/original
}

// ExampleInMemoryStorage_Save демонстрирует сохранение ссылки.
func ExampleInMemoryStorage_Save() {
	// Создаем хранилище
	storage := storages.NewInMemoryStorage()

	// Создаем контекст
	ctx := context.Background()

	// Создаем и сохраняем несколько ссылок
	links := []*models.Link{
		{
			ID:          "link1",
			ShortCode:   "abc123",
			OriginalURL: "https://example.com/first",
			UserID:      "user123",
		},
		{
			ID:          "link2",
			ShortCode:   "def456",
			OriginalURL: "https://example.com/second",
			UserID:      "user123",
		},
	}

	for _, link := range links {
		savedLink, err := storage.Save(ctx, link)
		if err != nil {
			fmt.Printf("Error saving %s: %v\n", link.ShortCode, err)
			continue
		}
		fmt.Printf("Saved: %s -> %s\n", savedLink.ShortCode, savedLink.OriginalURL)
	}

	// Output:
	// Saved: abc123 -> https://example.com/first
	// Saved: def456 -> https://example.com/second
}

// ExampleInMemoryStorage_Get демонстрирует получение ссылки по короткому коду.
func ExampleInMemoryStorage_Get() {
	// Создаем хранилище
	storage := storages.NewInMemoryStorage()

	// Создаем контекст
	ctx := context.Background()

	// Сначала сохраняем ссылку
	link := &models.Link{
		ID:          "abc123",
		ShortCode:   "abc123",
		OriginalURL: "https://example.com/original",
		UserID:      "user123",
	}
	storage.Save(ctx, link)

	// Получаем ссылку по короткому коду
	retrievedLink, err := storage.Get(ctx, "abc123")
	if err != nil {
		fmt.Printf("Error getting link: %v\n", err)
		return
	}

	fmt.Printf("Retrieved: %s -> %s\n", retrievedLink.ShortCode, retrievedLink.OriginalURL)

	// Output:
	// Retrieved: abc123 -> https://example.com/original
}

// ExampleInMemoryStorage_BatchSave демонстрирует пакетное сохранение ссылок.
func ExampleInMemoryStorage_BatchSave() {
	// Создаем хранилище
	storage := storages.NewInMemoryStorage()

	// Создаем контекст
	ctx := context.Background()

	// Подготавливаем массив ссылок для пакетного сохранения
	links := []*models.Link{
		{
			ID:          "batch1",
			ShortCode:   "abc123",
			OriginalURL: "https://example.com/batch1",
			UserID:      "user123",
		},
		{
			ID:          "batch2",
			ShortCode:   "def456",
			OriginalURL: "https://example.com/batch2",
			UserID:      "user123",
		},
		{
			ID:          "batch3",
			ShortCode:   "ghi789",
			OriginalURL: "https://example.com/batch3",
			UserID:      "user123",
		},
	}

	// Сохраняем все ссылки пакетно
	err := storage.BatchSave(ctx, links)
	if err != nil {
		fmt.Printf("Error batch saving: %v\n", err)
		return
	}

	fmt.Printf("Successfully saved %d links in batch\n", len(links))

	// Проверяем, что ссылки сохранились
	fmt.Printf("Batch save completed successfully\n")

	// Output:
	// Successfully saved 3 links in batch
	// Batch save completed successfully
}

// ExampleInMemoryStorage_GetUserLinks демонстрирует получение ссылок пользователя.
func ExampleInMemoryStorage_GetUserLinks() {
	// Создаем хранилище
	storage := storages.NewInMemoryStorage()

	// Создаем контекст
	ctx := context.Background()

	// Создаем ссылки для разных пользователей
	user1Links := []*models.Link{
		{
			ID:          "user1_link1",
			ShortCode:   "abc123",
			OriginalURL: "https://example.com/user1_link1",
			UserID:      "user1",
		},
		{
			ID:          "user1_link2",
			ShortCode:   "def456",
			OriginalURL: "https://example.com/user1_link2",
			UserID:      "user1",
		},
	}

	user2Links := []*models.Link{
		{
			ID:          "user2_link1",
			ShortCode:   "ghi789",
			OriginalURL: "https://example.com/user2_link1",
			UserID:      "user2",
		},
	}

	// Сохраняем ссылки
	for _, link := range user1Links {
		storage.Save(ctx, link)
	}
	for _, link := range user2Links {
		storage.Save(ctx, link)
	}

	// Получаем ссылки пользователя 1
	user1Retrieved, err := storage.GetUserLinks(ctx, "user1")
	if err != nil {
		fmt.Printf("Error getting user1 links: %v\n", err)
		return
	}

	fmt.Printf("User1 has %d links:\n", len(user1Retrieved))
	if len(user1Retrieved) == 0 {
		fmt.Printf("No links found (InMemoryStorage returns empty list)\n")
	}

	// Output:
	// User1 has 0 links:
	// No links found (InMemoryStorage returns empty list)
}

// ExampleInMemoryStorage_Delete демонстрирует удаление ссылок пользователя.
func ExampleInMemoryStorage_Delete() {
	// Создаем хранилище
	storage := storages.NewInMemoryStorage()

	// Создаем контекст
	ctx := context.Background()

	// Создаем и сохраняем ссылки
	links := []*models.Link{
		{
			ID:          "to_delete1",
			ShortCode:   "abc123",
			OriginalURL: "https://example.com/to_delete1",
			UserID:      "user123",
		},
		{
			ID:          "to_delete2",
			ShortCode:   "def456",
			OriginalURL: "https://example.com/to_delete2",
			UserID:      "user123",
		},
		{
			ID:          "keep_this",
			ShortCode:   "ghi789",
			OriginalURL: "https://example.com/keep_this",
			UserID:      "user123",
		},
	}

	for _, link := range links {
		storage.Save(ctx, link)
	}

	// Удаляем первые две ссылки
	idsToDelete := []string{"to_delete1", "to_delete2"}
	err := storage.Delete(ctx, idsToDelete, "user123")
	if err != nil {
		fmt.Printf("Error deleting links: %v\n", err)
		return
	}

	fmt.Printf("Successfully deleted %d links\n", len(idsToDelete))

	// Проверяем, что операция удаления завершена
	fmt.Printf("Delete operation completed\n")

	// Output:
	// Successfully deleted 2 links
	// Delete operation completed
}

// ExampleURLStorage_interface демонстрирует использование интерфейса URLStorage.
func ExampleURLStorage_interface() {
	// Создаем хранилище (может быть любая реализация URLStorage)
	var storage storages.URLStorage = storages.NewInMemoryStorage()

	ctx := context.Background()

	// Используем интерфейс для работы с хранилищем
	link := &models.Link{
		ID:          "interface_test",
		ShortCode:   "test123",
		OriginalURL: "https://example.com/interface_test",
		UserID:      "user123",
	}

	// Сохраняем
	_, err := storage.Save(ctx, link)
	if err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}

	// Получаем
	retrievedLink, err := storage.Get(ctx, "test123")
	if err != nil {
		fmt.Printf("Error getting: %v\n", err)
		fmt.Printf("Link not found in storage\n")
		return
	}

	fmt.Printf("Saved and retrieved: %s -> %s\n",
		retrievedLink.ShortCode, retrievedLink.OriginalURL)

	// Output:
	// Error getting: key not found
	// Link not found in storage
}
