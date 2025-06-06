package config

import (
	"flag"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type FlagProvider struct {
	HostFlagName            string
	ShortURLFlagName        string
	FileStoragePathFlagName string
	DatabaseDSNFlagName     string
}

func NewFlagProvider() *FlagProvider {
	return &FlagProvider{
		HostFlagName:            "a",
		ShortURLFlagName:        "b",
		FileStoragePathFlagName: "f",
		DatabaseDSNFlagName:     "d",
	}
}

func (flagConf *FlagProvider) setValues(c *Config) error {
	host := flag.String(flagConf.HostFlagName, "", "Адрес запуска HTTP-сервера")
	shortURLHost := flag.String(flagConf.ShortURLFlagName, "", "Базовый адрес результирующего сокращённого URL")
	fileStoragePath := flag.String(flagConf.FileStoragePathFlagName, "", "Путь к файлу для хранения")
	databaseDSN := flag.String(flagConf.DatabaseDSNFlagName, "", "Строка подлючения к БД")
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
