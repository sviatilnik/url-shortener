package config

import (
	"errors"
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
		return errors.New("invalid url")
	}

	return c.SetParam("urlBase", urlParam)
}

func (c ShortenerConfig) SetParam(key string, value params.ShortenerConfigParam) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key must not be blank")
	}

	c.params[key] = value
	return nil
}

func (c ShortenerConfig) GetParam(key string) (params.ShortenerConfigParam, error) {
	value, ok := c.params[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (c ShortenerConfig) GetParamValue(key string, default_value interface{}) interface{} {
	param, err := c.GetParam(key)
	if err != nil || !param.Validate() {
		return default_value
	}

	return param.Value()
}
