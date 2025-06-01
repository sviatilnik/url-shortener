package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"io"
	"net/http"
)

type batchRequestItem struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type batchResponseItem struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func BatchShortLinkHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

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

		links := make([]models.Link, len(req))
		for _, item := range req {
			links = append(links, models.Link{
				Id:          item.CorrelationId,
				OriginalURL: item.OriginalURL,
			})
		}

		generatedLinks, err := shorter.GenerateBatchShortLink(links)
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
				CorrelationId: item.Id,
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
		_, err = w.Write(encodedResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
