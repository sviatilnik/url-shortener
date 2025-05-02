package config

type DefaultProvider struct{}

func (d DefaultProvider) setValues(c *Config) error {
	c.Host = "localhost:8080"
	c.ShortURLHost = "http://localhost:8080"
	c.FileStoragePath = "store"
	return nil
}
