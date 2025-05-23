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
	shortenerConfig "github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"net/http"
	"time"
)

func main() {
	conf := getConfig()
	log := logger.NewLogger()

	connection, connErr := getDBConnection(&conf)
	if connErr != nil {
		log.Panicln("Failed to connect to database")
	}

	shorter := getShortener(conf.ShortURLHost, connection, &conf)

	r := chi.NewRouter()
	r.Use(middlewares.Log)
	r.Use(middlewares.Compress)
	r.Post("/", handlers.GetShortLinkHandler(shorter))
	r.Get("/{id}", handlers.RedirectToFullLinkHandler(shorter))
	r.Get("/ping", handlers.PingHandler(connection))
	r.Post("/api/shorten", handlers.APIShortLinkHandler(shorter))

	err := http.ListenAndServe(conf.Host, r)
	if err != nil {
		log.Fatalw("Error starting server", "error", err)
	}
}

func getShortener(baseURL string, db *sql.DB, config *config.Config) *shortener.Shortener {
	conf := shortenerConfig.NewShortenerConfig()
	_ = conf.SetURLBase(baseURL)

	return shortener.NewShortener(storages.NewPostgresStorageStorage(db), generators.NewRandomGenerator(10), conf)
}

func getConfig() config.Config {
	return config.NewConfig(&config.DefaultProvider{}, &config.FlagProvider{}, &config.EnvProvider{})
}

func getDBConnection(config *config.Config) (*sql.DB, error) {
	conn, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
