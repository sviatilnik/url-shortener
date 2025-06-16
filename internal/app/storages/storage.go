package storages

import (
	"context"
	"github.com/sviatilnik/url-shortener/internal/app/models"
)

type URLStorage interface {
	Save(ctx context.Context, link *models.Link) (*models.Link, error)
	BatchSave(ctx context.Context, links []*models.Link) error
	Get(ctx context.Context, shortCode string) (*models.Link, error)
	GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error)
	Delete(ctx context.Context, IDs []string, userID string) error
}
