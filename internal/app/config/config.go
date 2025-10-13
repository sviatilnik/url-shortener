package config

// Config представляет конфигурацию приложения.
// Содержит все необходимые параметры для работы сервиса сокращения URL.
type Config struct {
	Host            string // Адрес и порт для запуска HTTP-сервера
	ShortURLHost    string // Базовый URL для создания коротких ссылок
	FileStoragePath string // Путь к файлу для хранения данных (если используется файловое хранилище)
	DatabaseDSN     string // Строка подключения к базе данных
	AuthSecret      string // Секретный ключ для аутентификации
	AuditFile       string // Путь к файлу аудита
	AuditURL        string // URL для отправки аудита
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
