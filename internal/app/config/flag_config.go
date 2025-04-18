package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strconv"
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
	port := flag.String("a", "8000", "Адрес запуска HTTP-сервера")
	host := flag.String("b", "http://localhost:8000", "Базовый адрес результирующего сокращённого URL")
	flag.Parse()

	if strings.TrimSpace(*port) == "" {
		*port = "8000"
	}
	if _, err := strconv.Atoi(*port); err != nil {
		*port = "8000"
	}
	*port = strings.TrimLeft(*port, ":")

	if strings.TrimSpace(*host) == "" || !util.IsURL(*host) {
		*host = "http://localhost:8000"
	}

	_ = c.Set("port", *port)
	_ = c.Set("host", *host)
}
