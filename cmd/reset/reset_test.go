package main

import (
	"testing"
)

func TestUser_Reset(t *testing.T) {
	// Создаем пользователя с данными
	user := &User{
		ID:     123,
		Name:   "John Doe",
		Email:  "john@example.com",
		Active: true,
		Tags:   []string{"admin", "user"},
		Settings: map[string]string{
			"theme": "dark",
			"lang":  "en",
		},
		Profile: &Profile{
			Bio:    "Software developer",
			Avatar: "avatar.jpg",
			SocialLinks: map[string]string{
				"github": "https://github.com/john",
			},
			Metadata: []byte("some metadata"),
		},
	}

	// Проверяем, что данные установлены
	if user.ID != 123 {
		t.Errorf("Expected ID 123, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", user.Name)
	}
	if len(user.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(user.Tags))
	}
	if len(user.Settings) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(user.Settings))
	}
	if user.Profile == nil {
		t.Error("Expected Profile to be set")
	}

	// Сбрасываем состояние
	user.Reset()

	// Проверяем, что все поля сброшены
	if user.ID != 0 {
		t.Errorf("Expected ID 0, got %d", user.ID)
	}
	if user.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", user.Name)
	}
	if user.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", user.Email)
	}
	if user.Active != false {
		t.Errorf("Expected Active false, got %v", user.Active)
	}
	if len(user.Tags) != 0 {
		t.Errorf("Expected empty Tags slice, got length %d", len(user.Tags))
	}
	if len(user.Settings) != 0 {
		t.Errorf("Expected empty Settings map, got length %d", len(user.Settings))
	}
	if user.Profile != nil {
		t.Error("Expected Profile to be nil after reset")
	}
}

func TestProfile_Reset(t *testing.T) {
	profile := &Profile{
		Bio:    "Software developer",
		Avatar: "avatar.jpg",
		SocialLinks: map[string]string{
			"github": "https://github.com/john",
		},
		Metadata: []byte("some metadata"),
	}

	// Проверяем, что данные установлены
	if profile.Bio != "Software developer" {
		t.Errorf("Expected Bio 'Software developer', got '%s'", profile.Bio)
	}
	if len(profile.SocialLinks) != 1 {
		t.Errorf("Expected 1 social link, got %d", len(profile.SocialLinks))
	}
	if len(profile.Metadata) != 13 {
		t.Errorf("Expected 13 bytes metadata, got %d", len(profile.Metadata))
	}

	// Сбрасываем состояние
	profile.Reset()

	// Проверяем, что все поля сброшены
	if profile.Bio != "" {
		t.Errorf("Expected empty Bio, got '%s'", profile.Bio)
	}
	if profile.Avatar != "" {
		t.Errorf("Expected empty Avatar, got '%s'", profile.Avatar)
	}
	if len(profile.SocialLinks) != 0 {
		t.Errorf("Expected empty SocialLinks map, got length %d", len(profile.SocialLinks))
	}
	if len(profile.Metadata) != 0 {
		t.Errorf("Expected empty Metadata slice, got length %d", len(profile.Metadata))
	}
}

func TestCache_Reset(t *testing.T) {
	cache := &Cache{
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Keys:    []string{"key1", "key2"},
		Size:    100,
		MaxSize: 1000,
		Enabled: true,
		Backend: &Backend{
			URL:     "http://example.com",
			Timeout: 30,
			Retries: 3,
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
			Options: []string{"option1", "option2"},
		},
	}

	// Проверяем, что данные установлены
	if len(cache.Data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(cache.Data))
	}
	if len(cache.Keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(cache.Keys))
	}
	if cache.Size != 100 {
		t.Errorf("Expected Size 100, got %d", cache.Size)
	}
	if cache.Enabled != true {
		t.Errorf("Expected Enabled true, got %v", cache.Enabled)
	}
	if cache.Backend == nil {
		t.Error("Expected Backend to be set")
	}

	// Сбрасываем состояние
	cache.Reset()

	// Проверяем, что все поля сброшены
	if len(cache.Data) != 0 {
		t.Errorf("Expected empty Data map, got length %d", len(cache.Data))
	}
	if len(cache.Keys) != 0 {
		t.Errorf("Expected empty Keys slice, got length %d", len(cache.Keys))
	}
	if cache.Size != 0 {
		t.Errorf("Expected Size 0, got %d", cache.Size)
	}
	if cache.MaxSize != 0 {
		t.Errorf("Expected MaxSize 0, got %d", cache.MaxSize)
	}
	if cache.Enabled != false {
		t.Errorf("Expected Enabled false, got %v", cache.Enabled)
	}
	if cache.Backend != nil {
		t.Error("Expected Backend to be nil after reset")
	}
}

func TestBackend_Reset(t *testing.T) {
	backend := &Backend{
		URL:     "http://example.com",
		Timeout: 30,
		Retries: 3,
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Options: []string{"option1", "option2"},
	}

	// Проверяем, что данные установлены
	if backend.URL != "http://example.com" {
		t.Errorf("Expected URL 'http://example.com', got '%s'", backend.URL)
	}
	if backend.Timeout != 30 {
		t.Errorf("Expected Timeout 30, got %d", backend.Timeout)
	}
	if len(backend.Headers) != 1 {
		t.Errorf("Expected 1 header, got %d", len(backend.Headers))
	}
	if len(backend.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(backend.Options))
	}

	// Сбрасываем состояние
	backend.Reset()

	// Проверяем, что все поля сброшены
	if backend.URL != "" {
		t.Errorf("Expected empty URL, got '%s'", backend.URL)
	}
	if backend.Timeout != 0 {
		t.Errorf("Expected Timeout 0, got %d", backend.Timeout)
	}
	if backend.Retries != 0 {
		t.Errorf("Expected Retries 0, got %d", backend.Retries)
	}
	if len(backend.Headers) != 0 {
		t.Errorf("Expected empty Headers map, got length %d", len(backend.Headers))
	}
	if len(backend.Options) != 0 {
		t.Errorf("Expected empty Options slice, got length %d", len(backend.Options))
	}
}
