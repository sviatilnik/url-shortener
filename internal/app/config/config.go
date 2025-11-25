package config

// Config представляет конфигурацию приложения.
// Содержит все необходимые параметры для работы сервиса сокращения URL.
type Config struct {
	Host            string `env:"SERVER_ADDRESS"`      // Адрес и порт для запуска HTTP-сервера
	GRPCAddress     string `env:"GRPC_SERVER_ADDRESS"` // Адрес и порт для запуска gRPC-сервера
	ShortURLHost    string `env:"BASE_URL"`            // Базовый URL для создания коротких ссылок
	FileStoragePath string `env:"FILE_STORAGE_PATH"`   // Путь к файлу для хранения данных (если используется файловое хранилище)
	DatabaseDSN     string `env:"DATABASE_DSN"`        // Строка подключения к базе данных
	AuthSecret      string // Секретный ключ для аутентификации
	AuditFile       string `env:"AUDIT_FILE"`     // Путь к файлу аудита
	AuditURL        string `env:"AUDIT_URL"`      // URL для отправки аудита
	EnabledHTTPS    bool   `env:"ENABLE_HTTPS"`   // Сервер будет использовать SSL
	TrustedSubnet   string `env:"TRUSTED_SUBNET"` // CIDR подсеть для доступа к внутренним эндпоинтам
}

// NewConfig создает новую конфигурацию, объединяя значения из переданных провайдеров.
// Провайдеры применяются в порядке их передачи.
func NewConfig(providers ...Provider) Config {
	conf := Config{}
	for _, provider := range providers {
		_ = conf.setValues(provider)
	}

	return conf
}

func (c *Config) setValues(provider Provider) error {
	return provider.setValues(c)
}
