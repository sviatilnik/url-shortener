package storages

type URLStorage interface {
	Save(key string, value string) error
	Get(key string) (string, error)
}
