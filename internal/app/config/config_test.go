package config

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/config/mock_config"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		os.Args[0],
		"-a=https://google.com",
		"-d=database_dsn",
	}

	config := NewConfig(&DefaultProvider{}, &FlagProvider{}, NewEnvProvider(getMockEnvGetter(t)))

	assert.Equal(t, "https://google.com", config.Host)               // from flag provider
	assert.Equal(t, "https://short.google.com", config.ShortURLHost) // from env provider
	assert.Equal(t, "database_dsn", config.DatabaseDSN)              // from flag provider
	assert.Equal(t, "store", config.FileStoragePath)                 // from default provider
}

func getMockEnvGetter(t *testing.T) EnvGetter {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_config.NewMockEnvGetter(ctrl)
	m.EXPECT().LookupEnv("BASE_URL").Return("https://short.google.com", true).AnyTimes()
	m.EXPECT().LookupEnv(gomock.Any()).Return("", false).AnyTimes()

	return m
}
