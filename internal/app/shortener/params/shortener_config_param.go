package params

type ShortenerConfigParam interface {
	Validate() bool
	Value() interface{}
}
