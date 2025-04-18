package config

type Config interface {
	Get(key string, default_value interface{}) interface{}
	Set(key string, value interface{}) error
}
