package handlers

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"net/http"
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

		http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
	}
}
