package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"testing"
)

func TestShortener_GetFullLinkByID(t *testing.T) {

	i := storages.NewInMemoryStorage()
	err := i.Save("test", "http://google.com")
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		id        string
		want      string
		wantErr   bool
	}{
		{
			name:      "#1",
			storage:   i,
			generator: generators.NewRandomGenerator(10),
			id:        "test",
			want:      "http://google.com",
		},
		{
			name:      "#2",
			storage:   i,
			generator: generators.NewRandomGenerator(10),
			id:        "test3",
			wantErr:   true,
		},
		{
			name:      "#3",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			id:        "test",
			wantErr:   true,
		},
		{
			name:      "#4",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			id:        " ",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewShortener(tt.storage, tt.generator)
			got, err := shortener.GetFullLinkByID(tt.id)
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

	tests := []struct {
		name      string
		storage   storages.URLStorage
		generator generators.Generator
		url       string
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "#1",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			url:       "http://google.com",
			wantLen:   10,
		},
		{
			name:      "#2",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			url:       "test",
			wantErr:   true,
		},
		{
			name:      "#3",
			storage:   storages.NewInMemoryStorage(),
			generator: generators.NewRandomGenerator(10),
			url:       " ",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShortener(tt.storage, tt.generator)
			got, err := s.GenerateShortLink(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, tt.wantLen, len(got))
			}
		})
	}
}
