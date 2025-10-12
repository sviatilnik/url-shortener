package shortener

import (
	"context"
	"errors"
	"strings"

	"github.com/sviatilnik/url-shortener/internal/app/generators"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	"github.com/sviatilnik/url-shortener/internal/app/util"
)

type Shortener struct {
	storage   storages.URLStorage
	generator generators.Generator
	conf      Config
}

func NewShortener(store storages.URLStorage, generator generators.Generator, conf Config) *Shortener {
	return &Shortener{
		storage:   store,
		generator: generator,
		conf:      conf,
	}
}

func (s *Shortener) GetFullLinkByShortCode(ctx context.Context, shortCode string) (*models.Link, error) {
	if strings.TrimSpace(shortCode) == "" {
		return nil, ErrIDIsRequired
	}

	link, err := s.storage.Get(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Shortener) GenerateShortLink(ctx context.Context, url string) (string, error) {
	if !util.IsURL(url) {
		return "", ErrInvalidURL
	}
	userID := ""
	tmpUserID := ctx.Value(models.ContextUserID)
	if tmpUserID != nil {
		userID = tmpUserID.(string)
	}

	link := &models.Link{}
	var saveErr error
	var savedLink *models.Link
	short, err := s.generator.Get(url)
	if err != nil {
		return "", err
	}

	link.ID = short
	link.ShortCode = short
	link.OriginalURL = url
	link.UserID = userID

	savedLink, err = s.storage.Save(ctx, link)

	if errors.Is(err, storages.ErrOriginalURLAlreadyExists) {
		link.ShortCode = savedLink.ShortCode
		saveErr = ErrLinkConflict
	}

	if savedLink == nil {
		return "", ErrCreateShortLink
	}

	return s.getShortBase() + "/" + link.ShortCode, saveErr
}

func (s *Shortener) GenerateBatchShortLink(ctx context.Context, links []models.Link) ([]*models.Link, error) {
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

		link.ID = short
		link.ShortCode = short

		validLinks = append(validLinks, &link)
	}

	if len(validLinks) == 0 {
		return nil, ErrNoValidLinksInBatch
	}

	err := s.storage.BatchSave(ctx, validLinks)
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
	urlBase := s.conf.BaseURL
	return strings.TrimRight(urlBase, "/")
}

func (s *Shortener) GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error) {

	links, err := s.storage.GetUserLinks(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		link.ShortURL = s.getShortBase() + "/" + link.ShortCode
	}

	return links, nil
}

func (s *Shortener) DeleteUserLinks(ctx context.Context, linksIDs []string, userID string) error {
	return s.storage.Delete(ctx, linksIDs, userID)
}
