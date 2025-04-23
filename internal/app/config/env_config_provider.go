package config

import (
	"os"
	"strings"
)

type EnvConfigProvider struct{}

func (envConf *EnvConfigProvider) setValues(c *Config) error {
	host, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok && strings.TrimSpace(host) != "" {
		c.Host = host
	}

	shortURLHost, ok := os.LookupEnv("BASE_URL")
	if ok && strings.TrimSpace(shortURLHost) != "" {
		c.ShortURLHost = shortURLHost
	}

	return nil
}
