package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/logger"
	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	shortenerConfig "github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"net/http"
)

func main() {
	conf := getConfig()
	shorter := getShortener(conf.ShortURLHost, &conf)
	log := logger.NewLogger()

	r := chi.NewRouter()
	r.Use(middlewares.Log)
	r.Use(middlewares.Compress)
	r.Post("/", handlers.GetShortLinkHandler(shorter))
	r.Get("/{id}", handlers.RedirectToFullLinkHandler(shorter))

	apiShortLink := handlers.APIShortLink{Shortener: shorter}
	r.Post("/api/shorten", apiShortLink.Handler())

	host := conf.Host

	err := http.ListenAndServe(host, r)
	if err != nil {
		log.Fatalw("Error starting server", "error", err)
	}
}

func getShortener(baseURL string, config *config.Config) *shortener.Shortener {
	conf := shortenerConfig.NewShortenerConfig()
	_ = conf.SetURLBase(baseURL)

	return shortener.NewShortener(storages.NewFileStorage(config.FileStoragePath), generators.NewRandomGenerator(10), conf)
}

func getConfig() config.Config {
	return config.NewConfig(&config.DefaultProvider{}, &config.FlagProvider{}, &config.EnvProvider{})
}
