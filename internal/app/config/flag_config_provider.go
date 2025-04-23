package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type FlagConfigProvider struct{}

func (flagConf *FlagConfigProvider) setValues(c *Config) error {
	host := flag.String("a", "", "Адрес запуска HTTP-сервера")
	shortURLHost := flag.String("b", "", "Базовый адрес результирующего сокращённого URL")
	flag.Parse()

	if strings.TrimSpace(*host) != "" {
		c.Host = *host
	}

	if strings.TrimSpace(*shortURLHost) != "" && util.IsURL(*shortURLHost) {
		c.ShortURLHost = *shortURLHost
	}

	return nil
}
