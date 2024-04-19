package persist

type UserRepository interface {
	Repository[UserRecord, string]

	AddRoom(string, string) error
	RemoveRoom(string, string) error
}

type UserRecord struct {
	Id    string
	Name  string
	Rooms map[string]string // map[roomId]Name
}
