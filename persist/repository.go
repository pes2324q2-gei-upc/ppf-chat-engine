package persist

type Repository[Resource any, Key any] interface {
	Exists(pk Key) (bool, error)
	Add(record Resource) error
	Remove(pk Key) error
	Get(pk Key) (*Resource, error)
	GetAll() ([]*Resource, error)
}
