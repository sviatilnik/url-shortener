package app

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

var storage = make(map[string]string)

func GetLinkByID(id string) (string, error) {
	url, ok := storage[id]
	if !ok {
		return "", errors.New("not Found")
	}
	return url, nil
}

func GenerateShortLink(url string) (string, error) {
	short := createShortLink(url)

	for {
		err := saveLink(short, url)
		if err == nil {
			break
		}
	}

	return short, nil
}

func createShortLink(url string) string {
	hash := md5.Sum([]byte(url + time.Now().String()))
	short := hex.EncodeToString(hash[:])

	return short[:10]
}

func saveLink(id, url string) error {
	_, ok := storage[id]
	if ok {
		return errors.New("already Exists")
	}

	storage[id] = url
	return nil
}
