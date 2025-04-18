package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type FlagConfig struct {
	Config
}

func NewFlagConfig(confg Config) Config {
	if confg == nil {
		confg = NewDefaultConfig(nil)
	}
	conf := &FlagConfig{
		confg,
	}
	conf.parseFlags()

	return conf
}

func (c *FlagConfig) parseFlags() {
	host := flag.String("a", "", "Адрес запуска HTTP-сервера")
	shortURLHost := flag.String("b", "", "Базовый адрес результирующего сокращённого URL")
	flag.Parse()

	if strings.TrimSpace(*host) != "" {
		_ = c.Set("host", *host)
	}

	if strings.TrimSpace(*shortURLHost) == "" || !util.IsURL(*shortURLHost) {
		_ = c.Set("shortURLHost", *shortURLHost)
	}
}
