package handlers

import (
	"encoding/json"
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

func GetShortLinkApiHandler(shortener *shortener.Shortener) http.HandlerFunc {
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

		shortLink, err := shortener.GenerateShortLink(req.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		encodedResp, err := json.Marshal(response{
			Result: shortLink,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(encodedResp)
	}
}
