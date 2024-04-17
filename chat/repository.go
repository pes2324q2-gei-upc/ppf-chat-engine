package chat

type Repository[T any, K any] interface {
	Add(T) error
	Remove(string) error

	Get(string) (*K, error)
	GetAll() ([]K, error)
	Clear() error
}
