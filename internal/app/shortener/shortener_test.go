package shortener

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"testing"
)

func TestShortener_GetFullLinkByID(t *testing.T) {

	i := storages.NewInMemoryStorage()
	err := i.Save(&models.Link{
		Id:          "test",
		ShortCode:   "test",
		OriginalURL: "http://google.com",
	})
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		conf      config.ShortenerConfig
		shortCode string
		want      string
		wantErr   bool
	}{
		{
			name:      "#1",
			storage:   i,
			generator: generators.NewRandomGenerator(10),
			conf:      config.NewShortenerConfig(),
			shortCode: "test",
			want:      "http://google.com",
		},
		{
			name:      "#2",
			storage:   i,
			generator: generators.NewRandomGenerator(10),
			conf:      config.NewShortenerConfig(),
			shortCode: "test3",
			wantErr:   true,
		},
		{
			name:      "#3",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			conf:      config.NewShortenerConfig(),
			shortCode: "test",
			wantErr:   true,
		},
		{
			name:      "#4",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			conf:      config.NewShortenerConfig(),
			shortCode: " ",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewShortener(tt.storage, tt.generator, tt.conf)
			got, err := shortener.GetFullLinkByShortCode(tt.shortCode)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestShortener_GenerateShortLink(t *testing.T) {

	baseURL := "http://google.com/"
	con := config.NewShortenerConfig()
	con.SetURLBase(baseURL)

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		conf      config.ShortenerConfig
		url       string
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "#1",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			conf:      con,
			url:       "http://google.com",
			wantLen:   10 + len(baseURL),
		},
		{
			name:      "#2",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			conf:      con,
			url:       "test",
			wantErr:   true,
		},
		{
			name:      "#3",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			conf:      con,
			url:       " ",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShortener(tt.storage, tt.generator, tt.conf)
			got, err := s.GenerateShortLink(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				fmt.Println(got)
				assert.Equal(t, tt.wantLen, len(got))
			}
		})
	}
}
