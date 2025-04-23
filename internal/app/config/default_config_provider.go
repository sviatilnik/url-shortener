package config

type DefaultConfigProvider struct{}

func (d DefaultConfigProvider) setValues(c *Config) error {
	c.Host = "localhost:8080"
	c.ShortURLHost = "http://localhost:8080"
	return nil
}
