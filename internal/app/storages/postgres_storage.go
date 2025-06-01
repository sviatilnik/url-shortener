package storages

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"log"
	"strings"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorageStorage(db *sql.DB) InitableStorage {
	postgresStorage := new(PostgresStorage)
	postgresStorage.db = db

	return postgresStorage
}

func (p *PostgresStorage) Save(link *models.Link) error {
	if strings.TrimSpace(link.Id) == "" {
		return ErrEmptyKey
	}

	_, err := p.Get(link.Id)
	if errors.Is(err, ErrKeyNotFound) {
		_, err = p.db.ExecContext(
			context.Background(),
			"INSERT INTO link (\"uuid\", \"originalURL\", \"shortCode\") VALUES ($1, $2, $3)",
			link.Id, link.OriginalURL, link.ShortCode)
		return err
	}

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(context.Background(), "UPDATE link SET \"uuid\"=$1, \"originalURL\"=$2, \"shortCode\"=$3", link.Id, link.OriginalURL, link.ShortCode)
	return err
}

func (p *PostgresStorage) BatchSave(links []*models.Link) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, link := range links {
		_, err = tx.ExecContext(
			context.Background(),
			"INSERT INTO link (\"uuid\", \"originalURL\", \"shortCode\") VALUES ($1, $2, $3)",
			link.Id, link.OriginalURL, link.ShortCode)
		log.Println(err)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) Get(shortCode string) (*models.Link, error) {
	var row struct {
		uuid        string
		originalURL string
		shortCode   string
	}

	err := p.db.QueryRowContext(
		context.Background(),
		"SELECT * FROM link WHERE \"shortCode\"=$1 ORDER BY \"createdAt\" DESC LIMIT 1",
		shortCode).Scan(&row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyNotFound
	}

	if err != nil {
		return nil, err
	}

	return &models.Link{
		Id:          row.uuid,
		OriginalURL: row.originalURL,
		ShortCode:   row.shortCode,
	}, nil
}

func (p *PostgresStorage) Init() error {
	_, err := p.db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS link(
    "uuid" character varying(255) NOT NULL,
    "originalURL" character varying(512) NOT NULL,
    "shortCode" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("uuid"));
    CREATE INDEX IF NOT EXISTS "idx_link_shortCode" ON link ("shortCode");
`)
	return err
}
