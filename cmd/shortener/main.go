package main

import (
	"fmt"
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
	shorter := getShortener(conf.Get("host", "http://localhost:8000/").(string))

	r := chi.NewRouter()
	r.Post("/", GetShortLinkHandler(shorter))
	r.Get("/{id}", RedirectToFullLinkHandler(shorter))

	port := conf.Get("port", "8080").(string)
	fmt.Println("Listening on port " + port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func getShortener(baseUrl string) *shortener.Shortener {
	conf := shortenerConfig.NewShortenerConfig()
	_ = conf.SetURLBase(baseUrl)

	return shortener.NewShortener(storages.NewInMemoryStorage(), generators.NewRandomGenerator(10), conf)
}

func getConfig() config.Config {
	return config.NewFlagConfig()
}
