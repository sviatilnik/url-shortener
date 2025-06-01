package shortener

import "errors"

var (
	ErrInvalidURL       = errors.New("invalid url")
	ErrCreateShortLink  = errors.New("could not generate short link")
	ErrIDIsRequired     = errors.New("id is required")
	NoValidLinksInBatch = errors.New("no valid links in batch")
	NoLinksInBatch      = errors.New("no links in batch")
)
