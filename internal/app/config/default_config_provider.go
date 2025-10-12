package config

import (
	"crypto/rand"
	"encoding/hex"
)

type DefaultProvider struct{}

func (d *DefaultProvider) setValues(c *Config) error {
	c.Host = "localhost:8080"
	c.ShortURLHost = "http://localhost:8080"
	c.FileStoragePath = ""
	c.DatabaseDSN = ""
	c.AuthSecret = d.getAuthSecret()
	c.AuditFile = "audit.log"
	c.AuditURL = ""
	return nil
}

func (d *DefaultProvider) getAuthSecret() string {
	randBytes := make([]byte, 32)
	_, _ = rand.Read(randBytes)

	return hex.EncodeToString(randBytes)
}
