package handlers

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"net/http"
)

func RedirectToFullLinkHandler(shortener *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")

		fullLink, err := shortener.GetFullLinkByID(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
	}
}
