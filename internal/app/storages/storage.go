package storages

import (
	"context"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

// URLStorage определяет интерфейс для хранения ссылок.
// Интерфейс поддерживает операции создания, получения, пакетного сохранения и удаления ссылок.
type URLStorage interface {
	// Save сохраняет ссылку в хранилище.
	// Возвращает сохраненную ссылку или ошибку.
	Save(ctx context.Context, link *models.Link) (*models.Link, error)

	// BatchSave сохраняет массив ссылок в хранилище.
	// Операция выполняется атомарно - либо все ссылки сохраняются, либо ни одна.
	BatchSave(ctx context.Context, links []*models.Link) error

	// Get получает ссылку по короткому коду.
	// Возвращает ссылку или ошибку, если ссылка не найдена.
	Get(ctx context.Context, shortCode string) (*models.Link, error)

	// GetUserLinks получает все ссылки пользователя по его идентификатору.
	// Возвращает массив ссылок пользователя.
	GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error)

	// Delete помечает указанные ссылки как удаленные (soft delete).
	// Удаление выполняется только для ссылок, принадлежащих указанному пользователю.
	Delete(ctx context.Context, IDs []string, userID string) error
}
