package handlers

import (
	"context"
	"net/http"

	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
)

func RedirectToFullLinkHandler(shortener *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortCode := r.PathValue("short_code")

		link, err := shortener.GetFullLinkByShortCode(r.Context(), shortCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if link.IsDeleted {
			w.WriteHeader(http.StatusGone)
			return
		}

		// Устанавливаем URL в контекст для аудита
		ctx := context.WithValue(r.Context(), middlewares.AuditURLKey, link.OriginalURL)
		*r = *r.WithContext(ctx)

		http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
	}
}
