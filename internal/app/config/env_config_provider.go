package config

import (
	"os"
	"strings"
)

type EnvProvider struct{}

func (envConf *EnvProvider) setValues(c *Config) error {
	host, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok && strings.TrimSpace(host) != "" {
		c.Host = host
	}

	shortURLHost, ok := os.LookupEnv("BASE_URL")
	if ok && strings.TrimSpace(shortURLHost) != "" {
		c.ShortURLHost = shortURLHost
	}

	fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok && strings.TrimSpace(fileStoragePath) != "" {
		c.FileStoragePath = fileStoragePath
	}

	databaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok && strings.TrimSpace(databaseDSN) != "" {
		c.DatabaseDSN = databaseDSN
	}

	return nil
}
