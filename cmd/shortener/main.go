package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	shortenerConfig "github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"net/http"
)

func main() {
	conf := getConfig()
	shorter := getShortener(conf.Get("shortURLHost").(string))

	r := chi.NewRouter()
	r.Post("/", handlers.GetShortLinkHandler(shorter))
	r.Get("/{id}", handlers.RedirectToFullLinkHandler(shorter))

	host := conf.Get("host").(string)

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
	return config.NewEnvConfig(config.NewFlagConfig(nil))
}
