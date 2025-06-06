package storages

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"testing"
)

func TestPostgresStorage_BatchSave(t *testing.T) {

	//storage := NewPostgresStorageStorage(nil, "links_test")
	//tests := []struct {
	//	name     string
	//	links    []*models.Link
	//	wantErr  bool
	//}{
	//	{
	//		name:     "#1",
	//		links: []*models.Link{
	//			{
	//				ID:          "1",
	//				ShortCode:   "short_code1",
	//				OriginalURL: "original_url1",
	//			},
	//			{
	//				ID:          "2",
	//				ShortCode:   "short_code2",
	//				OriginalURL: "original_url2",
	//			},
	//		},
	//		wantErr: false,
	//	},
	//	{
	//		name:     "#2",
	//		links:    nil,
	//		wantErr:  true,
	//	},
	//	{
	//		name:     "#3",
	//		links: []*models.Link{
	//			{
	//				ID:          "1",
	//				ShortCode:   "short_code1",
	//				OriginalURL: "original_url1",
	//			},
	//			{
	//				ID:          "2",
	//				ShortCode:   "short_code2",
	//				OriginalURL: "original_url2",
	//			},
	//		},
	//		wantErr: true,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//
	//
	//		err := f.BatchSave(t.Context(), tt.links)
	//		if tt.wantErr {
	//			assert.Error(t, err)
	//		} else {
	//			assert.NoError(t, err)
	//		}
	//	})
	//}
	//
	//t.Cleanup(func() {
	//
	//})
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
