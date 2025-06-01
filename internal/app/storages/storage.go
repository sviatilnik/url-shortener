package storages

import "github.com/sviatilnik/url-shortener/internal/app/models"

type URLStorage interface {
	Save(link *models.Link) error
	BatchSave(links []*models.Link) error
	Get(shortCode string) (*models.Link, error)
}
