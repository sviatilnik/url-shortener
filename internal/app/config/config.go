package config

type Config interface {
	Get(key string) interface{}
	Set(key string, value interface{}) error
}
