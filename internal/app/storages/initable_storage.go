package storages

import "context"

type InitableStorage interface {
	URLStorage
	Init(ctx context.Context) error
}
