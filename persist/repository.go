package persist

// Repository represents a generic repository which can manage one or more database tables.
type Repository interface {
	Add(Record) error
	Remove(string) error

	Get(string) (*Record, error)
	GetAll() ([]Record, error)

	GetBy(string, string) (*Record, error) // GetBy returns a record by a specific field.
}
