package config

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/url-shortener/internal/app/config/mock_config"
)

func TestEnvProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_config.NewMockEnvGetter(ctrl)

	m.EXPECT().LookupEnv("SERVER_ADDRESS").Return("https://google.com", true).AnyTimes()
	m.EXPECT().LookupEnv("BASE_URL").Return("https://short.google.com", true).AnyTimes()
	m.EXPECT().LookupEnv("DATABASE_DSN").Return("database_dsn", true).AnyTimes()
	m.EXPECT().LookupEnv("FILE_STORAGE_PATH").Return("/tmp/file_storage", true).AnyTimes()
	m.EXPECT().LookupEnv("AUDIT_FILE").Return("audit-file", true).AnyTimes()
	m.EXPECT().LookupEnv("AUDIT_URL").Return("audit-url", true).AnyTimes()

	config := NewConfig(NewEnvProvider(m))

	assert.Equal(t, "https://google.com", config.Host)
	assert.Equal(t, "https://short.google.com", config.ShortURLHost)
	assert.Equal(t, "database_dsn", config.DatabaseDSN)
	assert.Equal(t, "/tmp/file_storage", config.FileStoragePath)
}
