package storages

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"sync"
	"testing"
)

func TestFileStorage_BatchSave(t *testing.T) {
	type fields struct {
		filePath string
		lastUUID int
		mut      sync.RWMutex
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
			f := &FileStorage{
				filePath: tt.fields.filePath,
				lastUUID: tt.fields.lastUUID,
				mut:      tt.fields.mut,
			}
			tt.wantErr(t, f.BatchSave(context.Background(), tt.args.links), fmt.Sprintf("BatchSave(%v)", tt.args.links))
		})
	}
}

func TestFileStorage_Get(t *testing.T) {
	type fields struct {
		filePath string
		lastUUID int
		mut      sync.RWMutex
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
			f := &FileStorage{
				filePath: tt.fields.filePath,
				lastUUID: tt.fields.lastUUID,
				mut:      tt.fields.mut,
			}
			got, err := f.Get(context.Background(), tt.args.shortCode)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.shortCode)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.shortCode)
		})
	}
}

func TestFileStorage_Save(t *testing.T) {
	type fields struct {
		filePath string
		lastUUID int
		mut      sync.RWMutex
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
			f := &FileStorage{
				filePath: tt.fields.filePath,
				lastUUID: tt.fields.lastUUID,
				mut:      tt.fields.mut,
			}
			got, err := f.Save(context.Background(), tt.args.link)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v)", tt.args.link)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Save(%v)", tt.args.link)
		})
	}
}

func TestNewFileStorage(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want *FileStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewFileStorage(tt.args.filePath), "NewFileStorage(%v)", tt.args.filePath)
		})
	}
}
