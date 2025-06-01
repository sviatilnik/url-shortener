package storages

type InitableStorage interface {
	URLStorage
	Init() error
}
