package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sviatilnik/url-shortener/internal/app/audit"
	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/handlers"
	"github.com/sviatilnik/url-shortener/internal/app/logger"
	"github.com/sviatilnik/url-shortener/internal/app/middlewares"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

	conf := getConfig()
	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	connection, connErr := getDBConnection(&conf)
	if connErr != nil {
		zapLogger.Info("Failed to connect to database")
	}

	storage := getStorage(ctx, connection, &conf)
	shorter := getShortener(conf.ShortURLHost, storage)
	auditService := getAuditService(&conf, zapLogger)

	r := chi.NewRouter()
	r.Use(middlewares.Log)
	r.Use(middlewares.Compress)
	r.Use(middlewares.NewAuthMiddleware(&conf, zapLogger).Auth)
	r.Use(middlewares.NewAuditMiddleware(auditService).Audit)

	if connection != nil {
		r.Get("/ping", handlers.PingDBHandler(connection))
	}
	r.Post("/", handlers.GetShortLinkHandler(shorter))
	r.Get("/{short_code}", handlers.RedirectToFullLinkHandler(shorter))
	r.Post("/api/shorten", handlers.APIShortLinkHandler(shorter))
	r.Post("/api/shorten/batch", handlers.BatchShortLinkHandler(shorter))
	r.Get("/api/user/urls", handlers.UserURLsHandler(shorter))
	r.Delete("/api/user/urls", handlers.DeleteUserURLsHandler(shorter))

	server := &http.Server{
		Addr:    conf.Host,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatalw("Error starting server", "error", err)
		}
	}()

	<-ctx.Done()

	zapLogger.Info("Shutting down server")

	timeout := 10 * time.Second
	ctxShutdown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		zapLogger.Fatalw("Error shutting down server", "error", err)
	}

	if connection != nil {
		if err := connection.Close(); err != nil {
			zapLogger.Fatalw("Error closing database connection", "error", err)
		}
		zapLogger.Info("Database connection closed successfully")
	}

	zapLogger.Info("Server shut down successfully")
}

func getShortener(baseURL string, storage storages.URLStorage) *shortener.Shortener {
	return shortener.NewShortener(
		storage,
		generators.NewRandomGenerator(10),
		shortener.NewShortenerConfig(baseURL),
	)
}

func getStorage(ctx context.Context, db *sql.DB, config *config.Config) storages.URLStorage {
	if db != nil {
		storage := storages.NewPostgresStorageStorage(db, "links")
		err := storage.Init(ctx)
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
	return config.NewConfig(
		&config.DefaultProvider{},
		config.NewFlagProvider(),
		config.NewEnvProvider(&config.OSEnvGetter{}),
	)
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

func getAuditService(config *config.Config, log *zap.SugaredLogger) *audit.AuditService {
	auditService := audit.NewAuditService(log)

	auditService.AddFileObserver(config.AuditFile)

	auditService.AddHTTPObserver(config.AuditURL)

	return auditService
}

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
