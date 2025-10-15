package models

// userID представляет тип для идентификатора пользователя в контексте.
type userID string

var (
	// ContextUserID используется как ключ для хранения идентификатора пользователя в контексте HTTP-запроса.
	ContextUserID userID
)
