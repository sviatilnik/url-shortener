package storages

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"testing"
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
			if tt.wantErr {
				assert.Error(t, i.Save(link))
			} else {
				assert.NoError(t, i.Save(link))
			}
		})
	}
}

func TestInMemoryStorage_Get(t *testing.T) {
	store := make(map[string]string)
	store["key"] = "value"
	store["key2"] = "value2"

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

			got, err := i.Get(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.OriginalURL)
			}
		})
	}
}
