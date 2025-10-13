package config_test

import (
	"fmt"

	"github.com/sviatilnik/url-shortener/internal/app/config"
)

// ExampleNewConfig демонстрирует создание конфигурации с несколькими провайдерами.
func ExampleNewConfig() {
	// Создаем провайдер по умолчанию
	defaultProvider := &config.DefaultProvider{}

	// Создаем провайдер переменных окружения
	envProvider := config.NewEnvProvider(&config.OSEnvGetter{})

	// Создаем провайдер флагов командной строки
	flagProvider := config.NewFlagProvider()

	// Создаем конфигурацию, объединяя все провайдеры
	conf := config.NewConfig(defaultProvider, envProvider, flagProvider)

	fmt.Printf("Host: %s\n", conf.Host)
	fmt.Printf("ShortURLHost: %s\n", conf.ShortURLHost)
	fmt.Printf("DatabaseDSN: %s\n", conf.DatabaseDSN)

	// Output:
	// Host: localhost:8080
	// ShortURLHost: http://localhost:8080
	// DatabaseDSN:
}

// ExampleConfig_priority демонстрирует приоритет провайдеров конфигурации.
func ExampleConfig_priority() {
	// Провайдеры применяются в порядке их передачи
	// Последний провайдер имеет наивысший приоритет

	// 1. Провайдер по умолчанию (низший приоритет)
	defaultProvider := &config.DefaultProvider{}

	// 2. Провайдер переменных окружения (средний приоритет)
	envProvider := config.NewEnvProvider(&config.OSEnvGetter{})

	// Создаем конфигурацию (без flagProvider чтобы избежать конфликта флагов)
	conf := config.NewConfig(defaultProvider, envProvider)

	fmt.Printf("Configuration created with %d providers\n", 2)
	fmt.Printf("Host from default: %s\n", conf.Host)

	// Output:
	// Configuration created with 2 providers
	// Host from default: localhost:8080
}

// ExampleConfig_fields демонстрирует различные поля конфигурации.
func ExampleConfig_fields() {
	// Создаем конфигурацию с провайдером по умолчанию
	conf := config.NewConfig(&config.DefaultProvider{})

	// Выводим основные поля конфигурации
	fmt.Printf("Server configuration:\n")
	fmt.Printf("- Host: %s\n", conf.Host)
	fmt.Printf("- ShortURLHost: %s\n", conf.ShortURLHost)
	fmt.Printf("- FileStoragePath: %s\n", conf.FileStoragePath)
	if conf.DatabaseDSN != "" {
		fmt.Printf("- DatabaseDSN: %s\n", conf.DatabaseDSN)
	} else {
		fmt.Printf("- DatabaseDSN:\n")
	}
	fmt.Printf("- AuthSecret length: %d\n", len(conf.AuthSecret))
	fmt.Printf("- AuditFile: %s\n", conf.AuditFile)
	fmt.Printf("- AuditURL: %s\n", conf.AuditURL)

	// Output:
	// Server configuration:
	// - Host: localhost:8080
	// - ShortURLHost: http://localhost:8080
	// - FileStoragePath: store
	// - DatabaseDSN:
	// - AuthSecret length: 64
	// - AuditFile: audit.log
	// - AuditURL:
}

// ExampleConfig_usage демонстрирует типичное использование конфигурации в приложении.
func ExampleConfig_usage() {
	// В реальном приложении конфигурация создается в main.go
	// и передается в различные компоненты

	conf := config.NewConfig(
		&config.DefaultProvider{},
		config.NewEnvProvider(&config.OSEnvGetter{}),
	)

	// Используем конфигурацию для настройки компонентов
	fmt.Printf("Starting server on %s\n", conf.Host)
	fmt.Printf("Short URLs will be created with base: %s\n", conf.ShortURLHost)

	if conf.DatabaseDSN != "" {
		fmt.Printf("Using database: %s\n", conf.DatabaseDSN)
	} else {
		fmt.Printf("Using in-memory storage\n")
	}

	if conf.FileStoragePath != "" {
		fmt.Printf("Using file storage: %s\n", conf.FileStoragePath)
	}

	// Output:
	// Starting server on localhost:8080
	// Short URLs will be created with base: http://localhost:8080
	// Using in-memory storage
	// Using file storage: store
}
