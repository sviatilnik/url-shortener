package storages

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type FileStorage struct {
	filePath string
	lastUUID int
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

func (f *FileStorage) Save(key string, value string) error {
	file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	item := &storeItem{
		OriginalURL: value,
		Short:       key,
		UUID:        key + "_" + strconv.Itoa(f.lastUUID),
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

func (f *FileStorage) Get(key string) (string, error) {
	file, err := os.Open(f.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		row = strings.Trim(row, "\n")

		item := &storeItem{}
		err = json.Unmarshal([]byte(row), item)
		if err != nil {
			return "", err
		}

		return item.OriginalURL, nil
	}

	return "", nil
}
