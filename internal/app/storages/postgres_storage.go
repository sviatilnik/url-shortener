package storages

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"strings"
)

type PostgresStorage struct {
	db        *sql.DB
	tableName string
}

func NewPostgresStorageStorage(db *sql.DB, tableName string) InitableStorage {
	postgresStorage := new(PostgresStorage)
	postgresStorage.db = db

	if strings.TrimSpace(tableName) != "" {
		tableName = strings.TrimSpace(tableName)
	} else {
		tableName = "links"
	}
	postgresStorage.tableName = tableName

	return postgresStorage
}

func (p *PostgresStorage) Save(ctx context.Context, link *models.Link) (*models.Link, error) {
	if strings.TrimSpace(link.ID) == "" {
		return nil, ErrEmptyKey
	}

	res, err := p.db.ExecContext(
		ctx,
		`INSERT INTO `+p.tableName+` ("uuid", "originalURL", "shortCode") 
				VALUES ($1, $2, $3) 
				ON CONFLICT("originalURL") DO NOTHING`,
		link.ID, link.OriginalURL, link.ShortCode)

	if err != nil {
		return nil, err
	}

	if c, _ := res.RowsAffected(); c != 1 {
		return p.GetByOriginalURL(ctx, link.OriginalURL), ErrOriginalURLAlreadyExists
	}

	return link, err
}

func (p *PostgresStorage) BatchSave(ctx context.Context, links []*models.Link) error {
	if len(links) == 0 {
		return ErrBatchIsEmpty
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, link := range links {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO `+p.tableName+` ("uuid", "originalURL", "shortCode") 
				    VALUES ($1, $2, $3) 
                    ON CONFLICT("uuid") DO UPDATE SET "originalURL" = $2, "shortCode" = $3`,
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

func (p *PostgresStorage) Get(ctx context.Context, shortCode string) (*models.Link, error) {
	var row struct {
		uuid        string
		originalURL string
		shortCode   string
	}

	err := p.db.QueryRowContext(
		ctx,
		`SELECT "uuid", "originalURL",  "shortCode" 
				FROM `+p.tableName+` 
				WHERE "shortCode"=$2`,
		p.tableName, shortCode).Scan(&row.uuid, &row.originalURL, &row.shortCode)
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

func (p *PostgresStorage) Init(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS `+p.tableName+` (
    "uuid" character varying(255) NOT NULL,
    "originalURL" character varying(512) NOT NULL,
    "shortCode" character varying(255) NOT NULL,
    "createdAt" timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("uuid"));
    CREATE INDEX IF NOT EXISTS "idx_link_shortCode" ON `+p.tableName+` ("shortCode");
	CREATE UNIQUE INDEX IF NOT EXISTS  "idx_link_originalUrl" ON `+p.tableName+` ("originalURL");`)
	return err
}

func (p *PostgresStorage) Drop(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, `DROP TABLE IF EXISTS `+p.tableName+`;`)
	return err
}

func (p *PostgresStorage) GetByOriginalURL(ctx context.Context, originalURL string) *models.Link {
	var row struct {
		uuid        string
		originalURL string
		shortCode   string
	}

	err := p.db.QueryRowContext(
		ctx,
		`SELECT "uuid", "originalURL",  "shortCode" 
				FROM `+p.tableName+` 
				WHERE "originalURL"=$1 
				ORDER BY "createdAt" DESC LIMIT 1`,
		originalURL).Scan(&row.uuid, &row.originalURL, &row.shortCode)

	if err != nil {
		return nil
	}

	return &models.Link{
		ID:          row.uuid,
		OriginalURL: row.originalURL,
		ShortCode:   row.shortCode,
	}
}
