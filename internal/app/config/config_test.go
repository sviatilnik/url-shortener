package config

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sviatilnik/url-shortener/internal/app/config/mock_config"
)

func TestConfig(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		os.Args[0],
		"-tta=https://google.com",
		"-ttd=database_dsn",
		"-ttauf=audit-file",
		"-ttau=audit-url",
	}

	config := NewConfig(
		&DefaultProvider{},
		&FlagProvider{
			HostFlagName:            "tta",
			ShortURLFlagName:        "ttb",
			FileStoragePathFlagName: "ttf",
			DatabaseDSNFlagName:     "ttd",
			AuditFileFlagName:       "ttauf",
			AuditURLFlagName:        "ttau",
		},
		NewEnvProvider(getMockEnvGetter(t)),
	)

	assert.Equal(t, "https://google.com", config.Host)               // from flag provider
	assert.Equal(t, "https://short.google.com", config.ShortURLHost) // from env provider
	assert.Equal(t, "database_dsn", config.DatabaseDSN)              // from flag provider
	assert.Equal(t, "store", config.FileStoragePath)                 // from default provider
	assert.Equal(t, "audit-file", config.AuditFile)                  // from flag provider
	assert.Equal(t, "audit-url", config.AuditURL)                    // from flag provider
}

func getMockEnvGetter(t *testing.T) EnvGetter {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_config.NewMockEnvGetter(ctrl)
	m.EXPECT().LookupEnv("BASE_URL").Return("https://short.google.com", true).AnyTimes()
	m.EXPECT().LookupEnv(gomock.Any()).Return("", false).AnyTimes()

	return m
}
