package config

import (
	"os"
	"strings"
)

type EnvConfig struct {
	Config
}

func NewEnvConfig(confg Config) Config {
	if confg == nil {
		confg = NewDefaultConfig(nil)
	}
	conf := &EnvConfig{
		confg,
	}
	conf.parseEnv()

	return conf
}

func (c *EnvConfig) parseEnv() {
	host, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok && strings.TrimSpace(host) != "" {
		_ = c.Set("host", host)
	}

	shortURLHost, ok := os.LookupEnv("BASE_URL")
	if ok && strings.TrimSpace(shortURLHost) != "" {
		_ = c.Set("shortURLHost", shortURLHost)
	}
}
