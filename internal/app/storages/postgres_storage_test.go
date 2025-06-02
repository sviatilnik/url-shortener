package storages

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"testing"
)

func TestNewPostgresStorageStorage(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want InitableStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPostgresStorageStorage(tt.args.db), "NewPostgresStorageStorage(%v)", tt.args.db)
		})
	}
}

func TestPostgresStorage_BatchSave(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		links []*models.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			tt.wantErr(t, p.BatchSave(context.Background(), tt.args.links), fmt.Sprintf("BatchSave(%v)", tt.args.links))
		})
	}
}

func TestPostgresStorage_Get(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		shortCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Link
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			got, err := p.Get(context.Background(), tt.args.shortCode)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.shortCode)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.shortCode)
		})
	}
}

func TestPostgresStorage_GetByOriginalURL(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		originalURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.Link
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			assert.Equalf(t, tt.want, p.GetByOriginalURL(context.Background(), tt.args.originalURL), "GetByOriginalURL(%v)", tt.args.originalURL)
		})
	}
}

func TestPostgresStorage_Init(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			tt.wantErr(t, p.Init(), fmt.Sprintf("Init()"))
		})
	}
}

func TestPostgresStorage_Save(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		link *models.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Link
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			got, err := p.Save(context.Background(), tt.args.link)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v)", tt.args.link)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Save(%v)", tt.args.link)
		})
	}
}

func TestPostgresStorage_isLinkExists(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				db: tt.fields.db,
			}
			assert.Equalf(t, tt.want, p.isLinkExists(tt.args.id), "isLinkExists(%v)", tt.args.id)
		})
	}
}
