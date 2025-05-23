package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type FlagProvider struct{}

func (flagConf *FlagProvider) setValues(c *Config) error {
	host := flag.String("a", "", "Адрес запуска HTTP-сервера")
	shortURLHost := flag.String("b", "", "Базовый адрес результирующего сокращённого URL")
	fileStoragePath := flag.String("f", "", "Путь к файлу для хранения")
	databaseDSN := flag.String("d", "", "Строка подлючения к БД")
	flag.Parse()

	if strings.TrimSpace(*host) != "" {
		c.Host = *host
	}

	if strings.TrimSpace(*shortURLHost) != "" && util.IsURL(*shortURLHost) {
		c.ShortURLHost = *shortURLHost
	}

	if strings.TrimSpace(*fileStoragePath) != "" {
		c.FileStoragePath = *fileStoragePath
	}

	if strings.TrimSpace(*databaseDSN) != "" {
		c.DatabaseDSN = *databaseDSN
	}

	return nil
}
