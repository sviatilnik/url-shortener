package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

// batchRequestItem представляет элемент запроса для пакетного создания коротких ссылок.
type batchRequestItem struct {
	CorrelationID string `json:"correlation_id"` // Идентификатор для связи запроса и ответа
	OriginalURL   string `json:"original_url"`   // Оригинальный URL для сокращения
}

// batchResponseItem представляет элемент ответа с созданной короткой ссылкой.
type batchResponseItem struct {
	CorrelationID string `json:"correlation_id"` // Идентификатор из запроса
	ShortURL      string `json:"short_url"`      // Сокращенная ссылка
}

// BatchShortLinkHandler создает HTTP-обработчик для пакетного создания коротких ссылок.
// Обработчик принимает массив JSON-объектов с полями "correlation_id" и "original_url"
// и возвращает массив JSON-объектов с полями "correlation_id" и "short_url".
// Возможные коды ответа:
//   - 201 Created - ссылки успешно созданы
//   - 400 Bad Request - неверный формат запроса или отсутствие валидных ссылок
//   - 500 Internal Server Error - внутренняя ошибка сервера
func BatchShortLinkHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req := make([]batchRequestItem, 0)
		err = json.Unmarshal(rawBody, &req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(models.ContextUserID).(string)

		links := make([]models.Link, len(req))
		for _, item := range req {
			links = append(links, models.Link{
				ID:          item.CorrelationID,
				OriginalURL: item.OriginalURL,
				UserID:      userID,
			})
		}

		generatedLinks, err := shorter.GenerateBatchShortLink(r.Context(), links)
		if errors.Is(err, shortener.ErrNoValidLinksInBatch) || errors.Is(err, shortener.ErrNoLinksInBatch) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := make([]batchResponseItem, 0)
		for _, item := range generatedLinks {
			resp = append(resp, batchResponseItem{
				CorrelationID: item.ID,
				ShortURL:      item.ShortURL,
			})
		}

		encodedResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(encodedResp)
	}
}
