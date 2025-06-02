package handlers

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"net/http"
)

func RedirectToFullLinkHandler(shortener *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortCode := r.PathValue("short_code")

		fullLink, err := shortener.GetFullLinkByShortCode(r.Context(), shortCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
	}
}
