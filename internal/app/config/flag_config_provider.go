package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/util"
)

type FlagProvider struct {
	HostFlagName            string
	GRPCAddressFlagName     string
	ShortURLFlagName        string
	FileStoragePathFlagName string
	DatabaseDSNFlagName     string
	AuditFileFlagName       string
	AuditURLFlagName        string
	EnablesHTTPSFlagName    string
	TrustedSubnetFlagName   string
}

func NewFlagProvider() *FlagProvider {
	return &FlagProvider{
		HostFlagName:            "a",
		GRPCAddressFlagName:     "g",
		ShortURLFlagName:        "b",
		FileStoragePathFlagName: "f",
		DatabaseDSNFlagName:     "d",
		AuditFileFlagName:       "audit-file",
		AuditURLFlagName:        "audit-url",
		EnablesHTTPSFlagName:    "s",
		TrustedSubnetFlagName:   "t",
	}
}

func (flagConf *FlagProvider) setValues(c *Config) error {
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	host := new(string)
	if name := strings.TrimSpace(flagConf.HostFlagName); name != "" {
		host = flagSet.String(name, "", "Адрес запуска HTTP-сервера")
	}

	grpcAddress := new(string)
	if name := strings.TrimSpace(flagConf.GRPCAddressFlagName); name != "" {
		grpcAddress = flagSet.String(name, "", "Адрес запуска gRPC-сервера")
	}

	shortURLHost := new(string)
	if name := strings.TrimSpace(flagConf.ShortURLFlagName); name != "" {
		shortURLHost = flagSet.String(name, "", "Базовый адрес результирующего сокращённого URL")
	}

	fileStoragePath := new(string)
	if name := strings.TrimSpace(flagConf.FileStoragePathFlagName); name != "" {
		fileStoragePath = flagSet.String(name, "", "Путь к файлу для хранения")
	}

	databaseDSN := new(string)
	if name := strings.TrimSpace(flagConf.DatabaseDSNFlagName); name != "" {
		databaseDSN = flagSet.String(name, "", "Строка подключения к БД")
	}

	auditFile := new(string)
	if name := strings.TrimSpace(flagConf.AuditFileFlagName); name != "" {
		auditFile = flagSet.String(name, "", "Путь к файлу для аудита")
	}

	auditURL := new(string)
	if name := strings.TrimSpace(flagConf.AuditURLFlagName); name != "" {
		auditURL = flagSet.String(name, "", "URL удаленного сервера для аудита")
	}

	enableHTTPS := new(bool)
	if name := strings.TrimSpace(flagConf.EnablesHTTPSFlagName); name != "" {
		enableHTTPS = flagSet.Bool(name, false, "Включить работу HTTPS")
	}

	trustedSubnet := new(string)
	if name := strings.TrimSpace(flagConf.TrustedSubnetFlagName); name != "" {
		trustedSubnet = flagSet.String(name, "", "CIDR подсеть для доступа к внутренним эндпоинтам")
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil && !errors.Is(err, flag.ErrHelp) {
		return err
	}

	if strings.TrimSpace(*host) != "" {
		c.Host = *host
	}

	if strings.TrimSpace(*grpcAddress) != "" {
		c.GRPCAddress = *grpcAddress
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

	if strings.TrimSpace(*auditFile) != "" {
		c.AuditFile = *auditFile
	}

	if strings.TrimSpace(*auditURL) != "" {
		c.AuditURL = *auditURL
	}

	c.EnabledHTTPS = *enableHTTPS

	if strings.TrimSpace(*trustedSubnet) != "" {
		c.TrustedSubnet = *trustedSubnet
	}

	return nil
}
