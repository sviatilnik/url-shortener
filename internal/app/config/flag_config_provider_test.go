package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFlagProvider(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		os.Args[0],
		"-a=https://google.com",
		"-b=https://short.google.com",
		"-f=/tmp/file_storage",
		"-d=database_dsn",
	}
	config := NewConfig(&FlagProvider{})

	assert.Equal(t, "https://google.com", config.Host)
	assert.Equal(t, "https://short.google.com", config.ShortURLHost)
	assert.Equal(t, "database_dsn", config.DatabaseDSN)
	assert.Equal(t, "/tmp/file_storage", config.FileStoragePath)
}
