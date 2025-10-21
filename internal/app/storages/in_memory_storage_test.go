package storages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

func TestInMemoryStorage_Save(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:  "#1",
			key:   "key",
			value: "value",
		},
		{
			name:    "#2",
			key:     " ",
			value:   "value",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInMemoryStorage()
			link := &models.Link{
				ID:        tt.key,
				ShortCode: tt.value,
			}
			_, err := i.Save(context.Background(), link)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInMemoryStorage_Get(t *testing.T) {
	store := make(map[string]*models.Link)
	store["key"] = &models.Link{
		ID:          "key",
		ShortCode:   "key",
		OriginalURL: "value",
		UserID:      "user1",
		IsDeleted:   false,
	}
	store["key2"] = &models.Link{
		ID:          "key2",
		ShortCode:   "key2",
		OriginalURL: "value2",
		UserID:      "user1",
		IsDeleted:   false,
	}

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name: "#1",
			key:  "key",
			want: "value",
		},
		{
			name:    "#2",
			key:     " ",
			wantErr: true,
		},
		{
			name:    "#3",
			key:     "key5",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := InMemoryStorage{
				store: store,
			}

			got, err := i.Get(context.Background(), tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.OriginalURL)
			}
		})
	}
}
