package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"io"
	"net/http"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result"`
}

func APIShortLinkHandler(short *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req := new(request)
		err = json.Unmarshal(rawBody, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status := http.StatusCreated
		shortLink, err := short.GenerateShortLink(r.Context(), req.URL)
		if err != nil {
			if errors.Is(err, shortener.ErrLinkConflict) {
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
				if errors.Is(err, shortener.ErrInvalidURL) {
					status = http.StatusBadRequest
				}
				w.WriteHeader(status)
				return
			}

		}
		encodedResp, err := json.Marshal(response{
			Result: shortLink,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(encodedResp)
	}
}
