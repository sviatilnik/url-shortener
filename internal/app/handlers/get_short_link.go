package handlers

import (
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"io"
	"net/http"
	"strings"
)

func GetShortLinkHandler(shorter *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status := http.StatusCreated
		shortLink, err := shorter.GenerateShortLink(r.Context(), strings.TrimSuffix(string(url), "\n"))
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
