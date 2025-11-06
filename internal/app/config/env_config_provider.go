package config

import (
	"os"
	"strings"
)

type EnvGetter interface {
	LookupEnv(key string) (string, bool)
}

type OSEnvGetter struct{}

func (O *OSEnvGetter) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

type EnvProvider struct {
	getter EnvGetter
}

func NewEnvProvider(getter EnvGetter) *EnvProvider {
	return &EnvProvider{
		getter: getter,
	}
}

func (env *EnvProvider) setValues(c *Config) error {
	host, ok := env.getter.LookupEnv("SERVER_ADDRESS")
	if ok && strings.TrimSpace(host) != "" {
		c.Host = host
	}

	shortURLHost, ok := env.getter.LookupEnv("BASE_URL")
	if ok && strings.TrimSpace(shortURLHost) != "" {
		c.ShortURLHost = shortURLHost
	}

	fileStoragePath, ok := env.getter.LookupEnv("FILE_STORAGE_PATH")
	if ok && strings.TrimSpace(fileStoragePath) != "" {
		c.FileStoragePath = fileStoragePath
	}

	databaseDSN, ok := env.getter.LookupEnv("DATABASE_DSN")
	if ok && strings.TrimSpace(databaseDSN) != "" {
		c.DatabaseDSN = databaseDSN
	}

	auditFile, ok := env.getter.LookupEnv("AUDIT_FILE")
	if ok && strings.TrimSpace(auditFile) != "" {
		c.AuditFile = auditFile
	}

	auditURL, ok := env.getter.LookupEnv("AUDIT_URL")
	if ok && strings.TrimSpace(auditURL) != "" {
		c.AuditURL = auditURL
	}

	enabledHTTPS, ok := env.getter.LookupEnv("ENABLE_HTTPS")
	if ok && strings.TrimSpace(enabledHTTPS) != "" {
		c.EnabledHTTPS = enabledHTTPS == "true"
	}

	return nil
}
