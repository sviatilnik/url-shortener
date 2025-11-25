package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

// statsResponse представляет ответ со статистикой сервиса.
type statsResponse struct {
	URLs  int `json:"urls"`  // Количество сокращённых URL в сервисе
	Users int `json:"users"` // Количество пользователей в сервисе
}

// InternalStatsHandler создает HTTP-обработчик для получения статистики сервиса.
// Обработчик возвращает JSON-объект с полями "urls" и "users".
// Возможные коды ответа:
//   - 200 OK - статистика успешно получена
//   - 403 Forbidden - доступ запрещён (IP не входит в доверенную подсеть)
//   - 500 Internal Server Error - внутренняя ошибка сервера
func InternalStatsHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		urls, users, err := shorter.GetStats(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := statsResponse{
			URLs:  urls,
			Users: users,
		}

		encodedResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encodedResp)
	}
}

