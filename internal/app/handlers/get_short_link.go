package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

// GetShortLinkHandler создает HTTP-обработчик для создания коротких ссылок через простой POST-запрос.
// Обработчик принимает URL в теле запроса и возвращает сокращенную ссылку в виде текста.
// Возможные коды ответа:
//   - 201 Created - ссылка успешно создана
//   - 409 Conflict - ссылка уже существует
//   - 400 Bad Request - неверный формат URL
func GetShortLinkHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		urlStr := strings.TrimSuffix(string(url), "\n")

		// Устанавливаем URL в контекст для аудита
		ctx := context.WithValue(r.Context(), middlewares.AuditURLKey, urlStr)
		*r = *r.WithContext(ctx)

		status := http.StatusCreated
		shortLink, err := shorter.GenerateShortLink(r.Context(), urlStr)
		if err != nil {
			if errors.Is(err, shortener.ErrLinkConflict) {
				status = http.StatusConflict
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.WriteHeader(status)
		w.Write([]byte(shortLink))
	}
}
