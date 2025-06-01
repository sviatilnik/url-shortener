package shortener

import (
	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener/config"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"github.com/sviatilnik/url-shortener/internal/app/util"
	"strings"
)

type Shortener struct {
	storage   storages.URLStorage
	generator generators.Generator
	conf      config.ShortenerConfig
}

func NewShortener(store storages.URLStorage, generator generators.Generator, conf config.ShortenerConfig) *Shortener {
	return &Shortener{
		storage:   store,
		generator: generator,
		conf:      conf,
	}
}

func (s *Shortener) GetFullLinkByShortCode(shortCode string) (string, error) {
	if strings.TrimSpace(shortCode) == "" {
		return "", ErrIDIsRequired
	}

	link, err := s.storage.Get(shortCode)
	if err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}

func (s *Shortener) GenerateShortLink(url string) (string, error) {
	if !util.IsURL(url) {
		return "", ErrInvalidURL
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

		link := &models.Link{
			Id:          short,
			ShortCode:   short,
			OriginalURL: url,
		}

		err = s.storage.Save(link)
		if err == nil {
			break
		}

		attemptsLeft--
	}

	if attemptsLeft <= 0 {
		return "", ErrCreateShortLink
	}

	return s.getShortBase() + "/" + short, nil
}

func (s *Shortener) GenerateBatchShortLink(links []models.Link) ([]*models.Link, error) {
	validLinks := make([]*models.Link, 0)

	if len(links) == 0 {
		return nil, ErrNoLinksInBatch
	}

	for _, link := range links {
		if !util.IsURL(link.OriginalURL) {
			continue
		}

		short, err := s.generator.Get(link.OriginalURL)
		if err != nil {
			continue
		}

		link.ShortCode = short

		validLinks = append(validLinks, &link)
	}

	if len(validLinks) == 0 {
		return nil, ErrNoValidLinksInBatch
	}

	err := s.storage.BatchSave(validLinks)
	if err != nil {
		return nil, err
	}

	shortBase := s.getShortBase()
	for _, link := range validLinks {
		link.ShortURL = shortBase + "/" + link.ShortCode
	}

	return validLinks, nil
}

func (s *Shortener) getShortBase() string {
	urlBase := s.conf.GetParamValue("urlBase", "http://localhost/").(string)
	return strings.TrimRight(urlBase, "/")
}
