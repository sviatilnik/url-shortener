package handlers

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"io"
	"net/http"
)

func GetShortLinkHandler(shortener *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		url, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortLink, err := shortener.GenerateShortLink(string(url))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortLink))
	}
}
