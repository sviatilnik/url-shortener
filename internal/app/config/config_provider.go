package config

type ConfigProvider interface {
	setValues(c *Config) error
}
