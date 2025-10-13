package models

// Link представляет структуру ссылки в системе сокращения URL.
type Link struct {
	ID          string // Уникальный идентификатор ссылки
	ShortCode   string // Короткий код для доступа к ссылке
	ShortURL    string // Полная сокращенная ссылка
	OriginalURL string // Оригинальный URL
	UserID      string // Идентификатор пользователя-владельца ссылки
	IsDeleted   bool   // Флаг удаления ссылки (soft delete)
}
