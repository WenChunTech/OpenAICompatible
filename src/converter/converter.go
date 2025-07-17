package converter

type Converter[T any] interface {
	Convert() (T, error)
}
