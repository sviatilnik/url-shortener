package shortener

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"github.com/sviatilnik/url-shortener/internal/app/storages/mock_storages"
)

func TestShortener_GetFullLinkByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storages.NewMockURLStorage(ctrl)
	mockStorage.EXPECT().Get(context.Background(), "test").Return(&models.Link{
		ID:          "test",
		ShortCode:   "test",
		OriginalURL: "http://google.com",
	}, nil).AnyTimes()

	mockStorage.EXPECT().Get(context.Background(), gomock.Any()).Return(nil, storages.ErrKeyNotFound).AnyTimes()

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		conf      Config
		shortCode string
		want      string
		wantErr   bool
	}{
		{
			name:      "#1",
			storage:   mockStorage,
			generator: generators.NewRandomGenerator(10),
			conf:      NewShortenerConfig(""),
			shortCode: "test",
			want:      "http://google.com",
		},
		{
			name:      "#2",
			storage:   mockStorage,
			generator: generators.NewRandomGenerator(10),
			conf:      NewShortenerConfig(""),
			shortCode: "test3",
			wantErr:   true,
		},
		{
			name:      "#3",
			storage:   mockStorage,
			generator: generators.NewRandomGenerator(10),
			conf:      NewShortenerConfig(""),
			shortCode: "test15",
			wantErr:   true,
		},
		{
			name:      "#4",
			storage:   mockStorage,
			generator: generators.NewRandomGenerator(10),
			conf:      NewShortenerConfig(""),
			shortCode: " ",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewShortener(tt.storage, tt.generator, tt.conf)
			got, err := shortener.GetFullLinkByShortCode(context.Background(), tt.shortCode)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.OriginalURL)
			}
		})
	}
}

func TestShortener_GenerateShortLink(t *testing.T) {

	con := NewShortenerConfig("http://google.com/")

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		conf      Config
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
			wantLen:   10 + len("http://google.com/"),
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
			got, err := s.GenerateShortLink(context.Background(), tt.url)
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
