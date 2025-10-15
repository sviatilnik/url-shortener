package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

// request представляет структуру запроса для создания короткой ссылки.
type request struct {
	URL string `json:"url"` // Оригинальный URL для сокращения
}

// response представляет структуру ответа с созданной короткой ссылкой.
type response struct {
	Result string `json:"result"` // Сокращенная ссылка
}

// APIShortLinkHandler создает HTTP-обработчик для API создания коротких ссылок.
// Обработчик принимает JSON-запрос с полем "url" и возвращает JSON-ответ с полем "result".
// Возможные коды ответа:
//   - 201 Created - ссылка успешно создана
//   - 409 Conflict - ссылка уже существует
//   - 400 Bad Request - неверный формат запроса или URL
//   - 500 Internal Server Error - внутренняя ошибка сервера
func APIShortLinkHandler(short *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req := new(request)
		err = json.Unmarshal(rawBody, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Устанавливаем URL в контекст для аудита
		ctx := context.WithValue(r.Context(), middlewares.AuditURLKey, req.URL)
		*r = *r.WithContext(ctx)

		status := http.StatusCreated
		shortLink, err := short.GenerateShortLink(r.Context(), req.URL)
		if err != nil {
			if errors.Is(err, shortener.ErrLinkConflict) {
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
				if errors.Is(err, shortener.ErrInvalidURL) {
					status = http.StatusBadRequest
				}
				w.WriteHeader(status)
				return
			}

		}
		encodedResp, err := json.Marshal(response{
			Result: shortLink,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(encodedResp)
	}
}
