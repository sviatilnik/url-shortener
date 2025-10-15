package shortener

import (
	"context"
	"testing"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
)

func BenchmarkShortener_GenerateShortLink(b *testing.B) {
	storage := storages.NewInMemoryStorage()
	generator := generators.NewRandomGenerator(10)
	conf := NewShortenerConfig("http://localhost:8080")
	shortener := NewShortener(storage, generator, conf)

	ctx := context.Background()
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shortener.GenerateShortLink(ctx, url)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkShortener_GetFullLinkByShortCode(b *testing.B) {
	storage := storages.NewInMemoryStorage()
	generator := generators.NewRandomGenerator(10)
	conf := NewShortenerConfig("http://localhost:8080")
	shortener := NewShortener(storage, generator, conf)

	ctx := context.Background()
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	// Предварительно создаем короткую ссылку
	shortURL, err := shortener.GenerateShortLink(ctx, url)
	if err != nil {
		b.Fatal(err)
	}

	// Извлекаем shortCode из URL
	shortCode := shortURL[len("http://localhost:8080/"):]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shortener.GetFullLinkByShortCode(ctx, shortCode)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkShortener_GenerateBatchShortLink(b *testing.B) {
	storage := storages.NewInMemoryStorage()
	generator := generators.NewRandomGenerator(10)
	conf := NewShortenerConfig("http://localhost:8080")
	shortener := NewShortener(storage, generator, conf)

	ctx := context.Background()
	links := []models.Link{
		{OriginalURL: "https://example.com/url1"},
		{OriginalURL: "https://example.com/url2"},
		{OriginalURL: "https://example.com/url3"},
		{OriginalURL: "https://example.com/url4"},
		{OriginalURL: "https://example.com/url5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shortener.GenerateBatchShortLink(ctx, links)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkShortener_GenerateShortLink_Parallel(b *testing.B) {
	storage := storages.NewInMemoryStorage()
	generator := generators.NewRandomGenerator(10)
	conf := NewShortenerConfig("http://localhost:8080")
	shortener := NewShortener(storage, generator, conf)

	ctx := context.Background()
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := shortener.GenerateShortLink(ctx, url)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
