package shortener

import "errors"

var (
	ErrInvalidURL          = errors.New("invalid url")
	ErrCreateShortLink     = errors.New("could not generate short link")
	ErrIDIsRequired        = errors.New("id is required")
	ErrNoValidLinksInBatch = errors.New("no valid links in batch")
	ErrNoLinksInBatch      = errors.New("no links in batch")
	ErrLinkConflict        = errors.New("link conflict")
)
