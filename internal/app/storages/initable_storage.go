package storages

import "context"

// InitableStorage расширяет интерфейс URLStorage возможностью инициализации.
// Используется для хранилищ, которые требуют предварительной настройки (например, создание таблиц в БД).
type InitableStorage interface {
	URLStorage
	// Init инициализирует хранилище (создает таблицы, индексы и т.д.).
	Init(ctx context.Context) error
}
