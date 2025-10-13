package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

// PingDBHandler создает HTTP-обработчик для проверки состояния базы данных.
// Обработчик выполняет ping-запрос к базе данных с таймаутом 3 секунды.
// Возможные коды ответа:
//   - 200 OK - база данных доступна
//   - 405 Method Not Allowed - неверный HTTP-метод
//   - 500 Internal Server Error - база данных недоступна
func PingDBHandler(conn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := conn.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
