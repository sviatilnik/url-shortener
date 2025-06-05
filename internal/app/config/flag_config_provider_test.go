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
		"-ta=https://google.com",
		"-tb=https://short.google.com",
		"-tf=/tmp/file_storage",
		"-td=database_dsn",
	}
	config := NewConfig(&FlagProvider{
		HostFlagName:            "ta",
		ShortURLFlagName:        "tb",
		FileStoragePathFlagName: "tf",
		DatabaseDSNFlagName:     "td",
	})

	assert.Equal(t, "https://google.com", config.Host)
	assert.Equal(t, "https://short.google.com", config.ShortURLHost)
	assert.Equal(t, "database_dsn", config.DatabaseDSN)
	assert.Equal(t, "/tmp/file_storage", config.FileStoragePath)
}
