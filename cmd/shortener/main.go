package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	shortenerConfig "github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
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

func main() {
	conf := getConfig()
	shorter := getShortener(conf.Get("shortUrlHost", "http://localhost:8080/").(string))

	r := chi.NewRouter()
	r.Post("/", GetShortLinkHandler(shorter))
	r.Get("/{id}", RedirectToFullLinkHandler(shorter))

	host := conf.Get("host", "localhost:8080").(string)

	err := http.ListenAndServe(host, r)
	if err != nil {
		panic(err)
	}
}

func getShortener(baseURL string) *shortener.Shortener {
	conf := shortenerConfig.NewShortenerConfig()
	_ = conf.SetURLBase(baseURL)

	return shortener.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10), conf)
}

func getConfig() config.Config {
	return config.NewFlagConfig()
}
