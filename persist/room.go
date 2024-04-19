package persist

type RoomRepository interface {
	Repository[RoomRecord, string]

	AddUser(string, string) error
	RemoveUser(string, string) error
}

type RoomRecord struct {
	Id    string
	Name  string
	Rooms map[string]string // map[roomId]Name
}
