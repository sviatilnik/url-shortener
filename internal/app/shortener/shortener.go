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

// Shortener представляет основной сервис для сокращения URL.
// Сервис использует генератор для создания коротких кодов и хранилище для сохранения ссылок.
type Shortener struct {
	storage   storages.URLStorage  // Хранилище для ссылок
	generator generators.Generator // Генератор коротких кодов
	conf      Config               // Конфигурация сервиса
}

// NewShortener создает новый экземпляр сервиса сокращения URL.
// Принимает хранилище, генератор и конфигурацию.
func NewShortener(store storages.URLStorage, generator generators.Generator, conf Config) *Shortener {
	return &Shortener{
		storage:   store,
		generator: generator,
		conf:      conf,
	}
}

// GetFullLinkByShortCode получает полную информацию о ссылке по короткому коду.
// Возвращает структуру Link с оригинальным URL и метаданными.
// Возможные ошибки:
//   - ErrIDIsRequired - короткий код не указан
//   - ErrKeyNotFound - ссылка не найдена
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

// GenerateShortLink создает короткую ссылку для указанного URL.
// Возвращает полную сокращенную ссылку.
// Возможные ошибки:
//   - ErrInvalidURL - неверный формат URL
//   - ErrLinkConflict - ссылка уже существует
//   - ErrCreateShortLink - ошибка создания ссылки
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

// GenerateBatchShortLink создает короткие ссылки для массива URL.
// Возвращает массив созданных ссылок с заполненными полями ShortURL.
// Возможные ошибки:
//   - ErrNoLinksInBatch - пустой массив ссылок
//   - ErrNoValidLinksInBatch - нет валидных URL в массиве
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

		link.ShortCode = short

		if link.ID == "" {
			link.ID = short
		}

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

// GetUserLinks получает все ссылки пользователя по его идентификатору.
// Возвращает массив ссылок с заполненными полями ShortURL.
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

// DeleteUserLinks помечает указанные ссылки как удаленные (soft delete).
// Принимает массив идентификаторов ссылок и идентификатор пользователя.
func (s *Shortener) DeleteUserLinks(ctx context.Context, linksIDs []string, userID string) error {
	return s.storage.Delete(ctx, linksIDs, userID)
}
