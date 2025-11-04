package config

import (
	"encoding/json"
	"os"
	"strings"
)

// JSONConfigProvider читает конфигурацию из JSON файла.
// Имя файла должно быть передано через метод SetConfigFilePath.
type JSONConfigProvider struct {
	configFilePath string
}

// NewJSONConfigProvider создает новый JSONConfigProvider.
// Если configFilePath пустой, провайдер не будет устанавливать значения.
func NewJSONConfigProvider(configFilePath string) *JSONConfigProvider {
	return &JSONConfigProvider{
		configFilePath: strings.TrimSpace(configFilePath),
	}
}

func (j *JSONConfigProvider) setValues(c *Config) error {
	// Если путь к файлу не указан, просто возвращаемся без ошибки
	if j.configFilePath == "" {
		return nil
	}

	data, err := os.ReadFile(j.configFilePath)
	if err != nil {
		// Если файл не существует, это не ошибка - просто не устанавливаем значения
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var jsonConfig struct {
		Host            string `json:"server_address"`
		ShortURLHost    string `json:"base_url"`
		FileStoragePath string `json:"file_storage_path"`
		DatabaseDSN     string `json:"database_dsn"`
		AuthSecret      string `json:"auth_secret"`
		AuditFile       string `json:"audit_file"`
		AuditURL        string `json:"audit_url"`
		EnabledHTTPS    bool   `json:"enable_https"`
	}

	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return err
	}

	if strings.TrimSpace(jsonConfig.Host) != "" {
		c.Host = jsonConfig.Host
	}

	if strings.TrimSpace(jsonConfig.ShortURLHost) != "" {
		c.ShortURLHost = jsonConfig.ShortURLHost
	}

	if strings.TrimSpace(jsonConfig.FileStoragePath) != "" {
		c.FileStoragePath = jsonConfig.FileStoragePath
	}

	if strings.TrimSpace(jsonConfig.DatabaseDSN) != "" {
		c.DatabaseDSN = jsonConfig.DatabaseDSN
	}

	if strings.TrimSpace(jsonConfig.AuthSecret) != "" {
		c.AuthSecret = jsonConfig.AuthSecret
	}

	if strings.TrimSpace(jsonConfig.AuditFile) != "" {
		c.AuditFile = jsonConfig.AuditFile
	}

	if strings.TrimSpace(jsonConfig.AuditURL) != "" {
		c.AuditURL = jsonConfig.AuditURL
	}

	if jsonConfig.EnabledHTTPS {
		c.EnabledHTTPS = jsonConfig.EnabledHTTPS
	}

	return nil
}
