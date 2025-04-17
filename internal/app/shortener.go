package app

import (
	"errors"
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type Shortener struct {
	storage   storages.URLStorage
	generator generators.Generator
}

func NewShortener(store storages.URLStorage, generator generators.Generator) *Shortener {
	return &Shortener{
		storage:   store,
		generator: generator,
	}
}

func (s *Shortener) GetFullLinkByID(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", errors.New("id is required")
	}

	return s.storage.Get(id)
}

func (s *Shortener) GenerateShortLink(url string) (string, error) {
	if !util.IsURL(url) {
		return "", errors.New("invalid url")
	}

	var short string
	var err error

	attemptsLeft := 10
	for {
		if attemptsLeft <= 0 {
			break
		}
		short, err = s.generator.Get(url)
		if err != nil {
			return "", err
		}

		err = s.storage.Save(short, url)
		if err == nil {
			break
		}

		attemptsLeft--
	}

	if attemptsLeft <= 0 {
		return "", errors.New("could not generate short link")
	}

	return short, nil
}
