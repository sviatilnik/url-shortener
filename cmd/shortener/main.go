package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sviatilnik/url-shortener/internal/app"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"io"
	"net/http"
)

func GetShortLinkHandler(shortener *app.Shortener) http.HandlerFunc {
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

		shortID, err := shortener.GenerateShortLink(string(url))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://" + r.Host + "/" + shortID))
	}
}

func RedirectToFullLinkHandler(shortener *app.Shortener) http.HandlerFunc {
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

func main() {
	shortener := getShortener()

	r := chi.NewRouter()
	r.Post("/", GetShortLinkHandler(shortener))
	r.Get("/{id}", RedirectToFullLinkHandler(shortener))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func getShortener() *app.Shortener {
	return app.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10))
}
