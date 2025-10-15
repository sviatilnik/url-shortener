package generators

import (
	"testing"
)

func BenchmarkRandomGenerator_Get(b *testing.B) {
	generator := NewRandomGenerator(10)
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.Get(url)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashGenerator_Get(b *testing.B) {
	generator := NewHashGenerator(10)
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.Get(url)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRandomGenerator_Get_Parallel(b *testing.B) {
	generator := NewRandomGenerator(10)
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := generator.Get(url)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkHashGenerator_Get_Parallel(b *testing.B) {
	generator := NewHashGenerator(10)
	url := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := generator.Get(url)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
