package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type FlagConfig struct {
	StaticConfig
}

func NewFlagConfig() Config {
	flagConfig := &FlagConfig{
		NewStaticConfig(nil),
	}
	flagConfig.parseFlags()

	return flagConfig
}

func (c *FlagConfig) parseFlags() {
	host := flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера")
	shortURLHost := flag.String("b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")
	flag.Parse()

	if strings.TrimSpace(*host) == "" {
		*host = "localhost:8080"
	}

	if strings.TrimSpace(*shortURLHost) == "" || !util.IsURL(*shortURLHost) {
		*shortURLHost = "http://localhost:8080"
	}

	_ = c.Set("host", *host)
	_ = c.Set("shortURLHost", *shortURLHost)
}
