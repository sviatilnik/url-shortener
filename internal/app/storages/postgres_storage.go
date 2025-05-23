package storages

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorageStorage(db *sql.DB) URLStorage {
	postgresStorage := new(PostgresStorage)
	postgresStorage.db = db

	postgresStorage.init()

	return postgresStorage
}

func (p *PostgresStorage) Save(short string, original string) error {
	if strings.TrimSpace(short) == "" {
		return ErrEmptyKey
	}

	_, err := p.Get(short)
	if errors.Is(err, ErrKeyNotFound) {
		_, err = p.db.ExecContext(
			context.Background(),
			"INSERT INTO link (\"uuid\", \"originalURL\", \"shortURL\") VALUES ($1, $2, $3)",
			short, original, short)
		return err
	}

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(context.Background(), "UPDATE link SET \"uuid\"=$1, \"originalURL\"=$2, \"shortURL\"=$3", short, original, short)
	return err
}

func (p *PostgresStorage) Get(key string) (string, error) {
	var fullURL string
	err := p.db.QueryRowContext(context.Background(), "SELECT \"originalURL\" FROM link WHERE \"uuid\"=$1", key).Scan(&fullURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrKeyNotFound
	}

	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func (p PostgresStorage) init() {
	p.db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS link(\"uuid\" character varying(255) NOT NULL,\"originalURL\" character varying(512) NOT NULL,\"shortURL\" character varying(255) NOT NULL,\"createdAt\" timestamp with time zone NOT NULL DEFAULT NOW(), PRIMARY KEY (\"uuid\"))")
}
