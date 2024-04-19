package persist

// Record represents a generic row in a database.
type Record interface {
	Pk() any
}
