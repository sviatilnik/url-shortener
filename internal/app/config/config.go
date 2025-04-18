package config

type Config interface {
	Get(key string, defaultValue interface{}) interface{}
	Set(key string, value interface{}) error
}
