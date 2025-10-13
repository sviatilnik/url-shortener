package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

type userURLsResponseItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func UserURLsHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpUserID := r.Context().Value(models.ContextUserID)
		if tmpUserID == nil {
			w.WriteHeader(http.StatusUnauthorized)
		}

		userID := tmpUserID.(string)
		if userID == "" || strings.TrimSpace(userID) == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}

		userLinks, err := shorter.GetUserLinks(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(userLinks) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp := make([]userURLsResponseItem, 0)
		for _, item := range userLinks {
			resp = append(resp, userURLsResponseItem{
				ShortURL:    item.ShortURL,
				OriginalURL: item.OriginalURL,
			})
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

func DeleteUserURLsHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpUserID := r.Context().Value(models.ContextUserID)
		if tmpUserID == nil {
			w.WriteHeader(http.StatusUnauthorized)
		}

		userID := tmpUserID.(string)
		if userID == "" || strings.TrimSpace(userID) == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		IDs := make([]string, 0)
		err = json.Unmarshal(rawBody, &IDs)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = shorter.DeleteUserLinks(r.Context(), IDs, userID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
