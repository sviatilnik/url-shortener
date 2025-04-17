package generators

type Generator interface {
	Get(str string) (string, error)
}
