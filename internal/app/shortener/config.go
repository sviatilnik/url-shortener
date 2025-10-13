package shortener

import "github.com/sviatilnik/url-shortener/internal/app/util"

// Config представляет конфигурацию сервиса сокращения URL.
type Config struct {
	BaseURL string // Базовый URL для создания коротких ссылок
}

// NewShortenerConfig создает новую конфигурацию сервиса сокращения URL.
// Если переданный BaseURL невалиден, используется значение по умолчанию "http://localhost/".
func NewShortenerConfig(BaseURL string) Config {
	if !util.IsURL(BaseURL) {
		return Config{
			BaseURL: "http://localhost/",
		}
	}
	return Config{
		BaseURL: BaseURL,
	}
}
