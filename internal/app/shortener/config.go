package shortener

import "github.com/sviatilnik/url-shortener/internal/app/util"

type Config struct {
	BaseURL string
}

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
