package storages

import (
	"context"
	"testing"

	"github.com/sviatilnik/url-shortener/internal/app/models"
)

func BenchmarkInMemoryStorage_Save(b *testing.B) {
	storage := NewInMemoryStorage()
	ctx := context.Background()
	link := &models.Link{
		ID:          "test123",
		ShortCode:   "test123",
		OriginalURL: "https://example.com/very/long/url/that/needs/to/be/shortened",
		UserID:      "user123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		link.ID = "test123" + string(rune(i))
		link.ShortCode = "test123" + string(rune(i))
		_, err := storage.Save(ctx, link)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryStorage_Get(b *testing.B) {
	storage := NewInMemoryStorage()
	ctx := context.Background()

	// Предварительно сохраняем ссылку
	link := &models.Link{
		ID:          "test123",
		ShortCode:   "test123",
		OriginalURL: "https://example.com/very/long/url/that/needs/to/be/shortened",
		UserID:      "user123",
	}
	_, err := storage.Save(ctx, link)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Get(ctx, "test123")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryStorage_BatchSave(b *testing.B) {
	storage := NewInMemoryStorage()
	ctx := context.Background()

	links := make([]*models.Link, 100)
	for i := 0; i < 100; i++ {
		links[i] = &models.Link{
			ID:          "test" + string(rune(i)),
			ShortCode:   "test" + string(rune(i)),
			OriginalURL: "https://example.com/url" + string(rune(i)),
			UserID:      "user123",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Обновляем ID для каждого запуска
		for j, link := range links {
			link.ID = "test" + string(rune(i*100+j))
			link.ShortCode = "test" + string(rune(i*100+j))
		}
		err := storage.BatchSave(ctx, links)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryStorage_Save_Parallel(b *testing.B) {
	storage := NewInMemoryStorage()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			link := &models.Link{
				ID:          "test" + string(rune(i)),
				ShortCode:   "test" + string(rune(i)),
				OriginalURL: "https://example.com/very/long/url/that/needs/to/be/shortened",
				UserID:      "user123",
			}
			_, err := storage.Save(ctx, link)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}
