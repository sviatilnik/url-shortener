package storages

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"os"
	"testing"
)

func TestFileStorage_BatchSave(t *testing.T) {
	file, tmpCreateErr := os.CreateTemp("", "test_file_storage")
	assert.NoError(t, tmpCreateErr)

	tests := []struct {
		name     string
		filePath string
		links    []*models.Link
		wantErr  bool
	}{
		{
			name:     "#1",
			filePath: file.Name(),
			links: []*models.Link{
				{
					ID:          "1",
					ShortCode:   "short_code1",
					OriginalURL: "original_url1",
				},
				{
					ID:          "2",
					ShortCode:   "short_code2",
					OriginalURL: "original_url2",
				},
			},
			wantErr: false,
		},
		{
			name:     "#2",
			filePath: file.Name(),
			links:    nil,
			wantErr:  true,
		},
		{
			name:     "#3",
			filePath: "",
			links: []*models.Link{
				{
					ID:          "1",
					ShortCode:   "short_code1",
					OriginalURL: "original_url1",
				},
				{
					ID:          "2",
					ShortCode:   "short_code2",
					OriginalURL: "original_url2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.filePath)

			err := f.BatchSave(t.Context(), tt.links)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})
}

func TestFileStorage_Get(t *testing.T) {
	file, tmpCreateErr := os.CreateTemp("", "test_file_storage")
	assert.NoError(t, tmpCreateErr)

	_, writeErr := file.Write([]byte("{\"uuid\":\"1\",\"short\":\"short_code\",\"original_url\":\"original_url\"}"))
	assert.NoError(t, writeErr)

	tests := []struct {
		name      string
		filePath  string
		shortCode string
		want      *models.Link
		wantErr   bool
	}{
		{
			name:      "#1",
			filePath:  file.Name(),
			shortCode: "short_code",
			want: &models.Link{
				ID:          "1",
				ShortCode:   "short_code",
				OriginalURL: "original_url",
			},
			wantErr: false,
		},
		{
			name:      "#2",
			filePath:  "",
			shortCode: "short_code",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "#3",
			filePath:  file.Name(),
			shortCode: "",
			want:      nil,
			wantErr:   false,
		},
		{
			name:      "#4",
			filePath:  file.Name(),
			shortCode: "short_code2",
			want:      nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.filePath)
			got, err := f.Get(context.Background(), tt.shortCode)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})
}

func TestFileStorage_Save(t *testing.T) {
	file, tmpCreateErr := os.CreateTemp("", "test_file_storage")
	assert.NoError(t, tmpCreateErr)

	tests := []struct {
		name     string
		filePath string
		link     *models.Link
		wantErr  bool
	}{
		{
			name:     "#1",
			filePath: file.Name(),
			link: &models.Link{
				ID:          "1",
				ShortURL:    "",
				ShortCode:   "short_code",
				OriginalURL: "original_url",
			},
			wantErr: false,
		},
		{
			name:     "#2",
			filePath: "",
			link:     &models.Link{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.filePath)
			_, err := f.Save(context.Background(), tt.link)

			if tt.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})
}
