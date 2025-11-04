package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONConfigProvider(t *testing.T) {
	// Создаем временный файл конфигурации
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	// Тест 1: Чтение полной конфигурации из JSON
	t.Run("read full config from JSON", func(t *testing.T) {
		jsonConfig := `{
			"server_address": "localhost:9090",
			"base_url": "http://localhost:9090",
			"file_storage_path": "/tmp/test_storage",
			"database_dsn": "postgres://test",
			"auth_secret": "test-secret",
			"audit_file": "/tmp/audit.log",
			"audit_url": "http://audit.example.com",
			"enable_https": true
		}`

		err := os.WriteFile(configFile, []byte(jsonConfig), 0644)
		require.NoError(t, err)

		config := NewConfig(
			&DefaultProvider{},
			NewJSONConfigProvider(configFile),
		)

		assert.Equal(t, "localhost:9090", config.Host)
		assert.Equal(t, "http://localhost:9090", config.ShortURLHost)
		assert.Equal(t, "/tmp/test_storage", config.FileStoragePath)
		assert.Equal(t, "postgres://test", config.DatabaseDSN)
		assert.Equal(t, "test-secret", config.AuthSecret)
		assert.Equal(t, "/tmp/audit.log", config.AuditFile)
		assert.Equal(t, "http://audit.example.com", config.AuditURL)
		assert.Equal(t, true, config.EnabledHTTPS)
	})

	// Тест 2: Чтение частичной конфигурации из JSON
	t.Run("read partial config from JSON", func(t *testing.T) {
		jsonConfig := `{
			"server_address": "localhost:7070",
			"database_dsn": "postgres://partial"
		}`

		err := os.WriteFile(configFile, []byte(jsonConfig), 0644)
		require.NoError(t, err)

		config := NewConfig(
			&DefaultProvider{},
			NewJSONConfigProvider(configFile),
		)

		// Значения из JSON
		assert.Equal(t, "localhost:7070", config.Host)
		assert.Equal(t, "postgres://partial", config.DatabaseDSN)
		// Значения по умолчанию из DefaultProvider
		assert.Equal(t, "http://localhost:8080", config.ShortURLHost)
		assert.Equal(t, "store", config.FileStoragePath)
	})

	// Тест 3: Файл не существует - не должно быть ошибки
	t.Run("non-existent file", func(t *testing.T) {
		config := NewConfig(
			&DefaultProvider{},
			NewJSONConfigProvider("/non/existent/file.json"),
		)

		// Должны остаться значения по умолчанию
		assert.Equal(t, "localhost:8080", config.Host)
		assert.Equal(t, "http://localhost:8080", config.ShortURLHost)
	})

	// Тест 4: Пустой путь к файлу
	t.Run("empty file path", func(t *testing.T) {
		config := NewConfig(
			&DefaultProvider{},
			NewJSONConfigProvider(""),
		)

		// Должны остаться значения по умолчанию
		assert.Equal(t, "localhost:8080", config.Host)
		assert.Equal(t, "http://localhost:8080", config.ShortURLHost)
	})

	// Тест 5: Невалидный JSON
	t.Run("invalid JSON", func(t *testing.T) {
		invalidJSON := `{invalid json}`
		err := os.WriteFile(configFile, []byte(invalidJSON), 0644)
		require.NoError(t, err)

		config := Config{}
		provider := NewJSONConfigProvider(configFile)
		err = provider.setValues(&config)
		assert.Error(t, err)
	})

	// Тест 6: Приоритет JSON над DefaultProvider, но ниже FlagProvider
	t.Run("priority JSON over default but below flags", func(t *testing.T) {
		jsonConfig := `{
			"server_address": "localhost:9090",
			"base_url": "http://localhost:9090"
		}`

		err := os.WriteFile(configFile, []byte(jsonConfig), 0644)
		require.NoError(t, err)

		// Создаем конфигурацию с правильным порядком провайдеров
		// JSON должен перезаписать Default, но Flag должен перезаписать JSON
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()

		os.Args = []string{
			os.Args[0],
			"-a=localhost:8081",
		}

		config := NewConfig(
			&DefaultProvider{},
			NewJSONConfigProvider(configFile),
			NewFlagProvider(),
		)

		// Флаг должен иметь приоритет над JSON
		assert.Equal(t, "localhost:8081", config.Host)
		// JSON должен иметь приоритет над Default
		assert.Equal(t, "http://localhost:9090", config.ShortURLHost)
	})
}
