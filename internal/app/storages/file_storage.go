package storages

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/sviatilnik/url-shortener/internal/app/models"
	"os"
	"strings"
	"sync"
)

type FileStorage struct {
	filePath string
	lastUUID int
	mut      sync.RWMutex
}

type storeItem struct {
	UUID        string `json:"uuid"`
	Short       string `json:"short"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
		lastUUID: 0,
	}
}

func (f *FileStorage) Save(ctx context.Context, link *models.Link) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

		file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		f.mut.Lock()
		defer f.mut.Unlock()

		item := &storeItem{
			OriginalURL: link.OriginalURL,
			Short:       link.ShortCode,
			UUID:        link.ID,
		}

		marshal, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		if _, err = file.Write(marshal); err != nil {
			return nil, err
		}

		file.WriteString("\n")

		f.lastUUID++
		return link, nil
	}
}

func (f *FileStorage) BatchSave(ctx context.Context, links []*models.Link) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if len(links) == 0 {
			return ErrBatchIsEmpty
		}

		for _, link := range links {
			if _, err := f.Save(ctx, link); err != nil {
				return err
			}
		}

		return nil
	}
}

func (f *FileStorage) Get(ctx context.Context, shortCode string) (*models.Link, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		file, err := os.Open(f.filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			row = strings.Trim(row, "\n")

			item := &storeItem{}
			err = json.Unmarshal([]byte(row), item)
			if err != nil {
				return nil, err
			}

			if item.Short == shortCode {
				return &models.Link{
					ID:          item.UUID,
					ShortCode:   item.Short,
					OriginalURL: item.OriginalURL,
				}, nil
			}
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (f *FileStorage) GetUserLinks(ctx context.Context, userID string) ([]*models.Link, error) {
	panic("implement me")
}
