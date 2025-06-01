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
	if strings.TrimSpace(link.ID) == "" {
		return ErrEmptyKey
	}

	if !p.isLinkExists(link.ID) {
		_, err := p.db.ExecContext(
			context.Background(),
			"INSERT INTO link (\"uuid\", \"originalURL\", \"shortCode\") VALUES ($1, $2, $3)",
			link.ID, link.OriginalURL, link.ShortCode)
		return err
	}

	_, err := p.db.ExecContext(context.Background(),
		"UPDATE link SET \"originalURL\"=$1, \"shortCode\"=$2 WHERE \"uuid\"=$3",
		link.OriginalURL, link.ShortCode, link.ID)

	return err
}

func (p *PostgresStorage) BatchSave(links []*models.Link) error {
	tx, err := p.db.Begin()
	log.Println(err)
	if err != nil {
		return err
	}

	for _, link := range links {
		if p.isLinkExists(link.ID) {
			_, err = p.db.ExecContext(
				context.Background(),
				"UPDATE link SET \"originalURL\"=$1, \"shortCode\"=$2 WHERE \"uuid\"=$3",
				link.OriginalURL, link.ShortCode, link.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
			continue
		}
		_, err = tx.ExecContext(
			context.Background(),
			"INSERT INTO link (\"uuid\", \"originalURL\", \"shortCode\") VALUES ($1, $2, $3)",
			link.ID, link.OriginalURL, link.ShortCode)
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
		`SELECT "uuid", "originalURL",  "shortCode" 
				FROM link 
				WHERE "shortCode"=$1 
				ORDER BY "createdAt" DESC LIMIT 1`,
		shortCode).Scan(&row.uuid, &row.originalURL, &row.shortCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyNotFound
	}

	if err != nil {
		return nil, err
	}

	return &models.Link{
		ID:          row.uuid,
		OriginalURL: row.originalURL,
		ShortCode:   row.shortCode,
	}, nil
}

func (p *PostgresStorage) isLinkExists(id string) bool {
	var dummy interface{}
	row := p.db.QueryRowContext(context.Background(), "SELECT 1 FROM link WHERE \"uuid\"=$1 ", id)
	err := row.Scan(&dummy)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false
	case err != nil:
		return false
	}

	return true
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
