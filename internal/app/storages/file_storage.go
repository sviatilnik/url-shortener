package storages

import (
	"bufio"
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

func (f *FileStorage) Save(link *models.Link) error {
	file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	f.mut.Lock()
	defer f.mut.Unlock()

	item := &storeItem{
		OriginalURL: link.OriginalURL,
		Short:       link.ShortCode,
		UUID:        link.Id,
	}

	marshal, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if _, err := file.Write(marshal); err != nil {
		return err
	}

	file.WriteString("\n")

	f.lastUUID++
	return nil
}

func (f *FileStorage) BatchSave(links []*models.Link) error {
	for _, link := range links {
		if err := f.Save(link); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileStorage) Get(shortCode string) (*models.Link, error) {
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
				Id:          item.UUID,
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
