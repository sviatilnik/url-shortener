package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProvider(t *testing.T) {
	config := NewConfig(&DefaultProvider{})

	assert.Equal(t, "localhost:8080", config.Host)
	assert.Equal(t, "http://localhost:8080", config.ShortURLHost)
	assert.Equal(t, "", config.DatabaseDSN)
	assert.Equal(t, "store", config.FileStoragePath)
}
