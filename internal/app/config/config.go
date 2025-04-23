package config

type Config struct {
	Host         string
	ShortURLHost string
}

func NewConfig(providers ...ConfigProvider) Config {
	conf := Config{}
	for _, provider := range providers {
		_ = conf.setValues(provider)
	}

	return conf
}

func (c *Config) setValues(provider ConfigProvider) error {
	return provider.setValues(c)
}
