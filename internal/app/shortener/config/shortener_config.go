package config

import (
	"github.com/sviatilnik/url-shortener/internal/app/shortener/params"
	"strings"
)

type ShortenerConfig struct {
	params map[string]params.ShortenerConfigParam
}

func NewShortenerConfig() ShortenerConfig {
	return ShortenerConfig{
		params: make(map[string]params.ShortenerConfigParam),
	}
}

func (c ShortenerConfig) SetURLBase(urlBase string) error {
	urlParam := params.NewURLConfigParam(urlBase)
	if !urlParam.Validate() {
		return ErrInvalidURL
	}

	return c.SetParam("urlBase", urlParam)
}

func (c ShortenerConfig) SetParam(key string, value params.ShortenerConfigParam) error {
	if strings.TrimSpace(key) == "" {
		return ErrKeyIsEmpty
	}

	c.params[key] = value
	return nil
}

func (c ShortenerConfig) GetParam(key string) (params.ShortenerConfigParam, error) {
	value, ok := c.params[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

func (c ShortenerConfig) GetParamValue(key string, defaultValue interface{}) interface{} {
	param, err := c.GetParam(key)
	if err != nil || !param.Validate() {
		return defaultValue
	}

	return param.Value()
}
