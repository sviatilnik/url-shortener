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

type ApiShortLink struct {
	Shortener *shortener.Shortener
}

func (h *ApiShortLink) Handler() http.HandlerFunc {
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
		req := new(request)
		err = json.Unmarshal(rawBody, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortLink, err := h.Shortener.GenerateShortLink(req.URL)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, shortener.ErrInvalidURL) {
				status = http.StatusBadRequest
			}

			w.WriteHeader(status)
			return
		}
		encodedResp, err := json.Marshal(response{
			Result: shortLink,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(encodedResp)
	}
}
