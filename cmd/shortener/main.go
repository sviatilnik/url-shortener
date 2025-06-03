package main

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/logger"
	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"net/http"
	"time"
)

func main() {
	conf := getConfig()
	log := logger.NewLogger()

	connection, connErr := getDBConnection(&conf)
	if connErr != nil {
		log.Info("Failed to connect to database")
	}

	storage := getStorage(connection, &conf)
	shorter := getShortener(conf.ShortURLHost, storage)

	r := chi.NewRouter()
	r.Use(middlewares.Log)
	r.Use(middlewares.Compress)
	if connection != nil {
		r.Get("/ping", handlers.PingDBHandler(connection))
	}
	r.Post("/", handlers.GetShortLinkHandler(shorter))
	r.Get("/{short_code}", handlers.RedirectToFullLinkHandler(shorter))
	r.Post("/api/shorten", handlers.APIShortLinkHandler(shorter))
	r.Post("/api/shorten/batch", handlers.BatchShortLinkHandler(shorter))

	err := http.ListenAndServe(conf.Host, r)
	if err != nil {
		log.Fatalw("Error starting server", "error", err)
	}
}

func getShortener(baseURL string, storage storages.URLStorage) *shortener.Shortener {
	return shortener.NewShortener(
		storage,
		generators.NewRandomGenerator(10),
		shortener.NewShortenerConfig(baseURL),
	)
}

func getStorage(db *sql.DB, config *config.Config) storages.URLStorage {
	if db != nil {
		storage := storages.NewPostgresStorageStorage(db)
		err := storage.Init()
		if err != nil {
			return nil
		}

		return storage
	}

	if config.FileStoragePath != "" {
		return storages.NewFileStorage(config.FileStoragePath)
	}

	return storages.NewInMemoryStorage()
}

func getConfig() config.Config {
	return config.NewConfig(&config.DefaultProvider{}, &config.FlagProvider{}, &config.EnvProvider{})
}

func getDBConnection(config *config.Config) (*sql.DB, error) {
	conn, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
