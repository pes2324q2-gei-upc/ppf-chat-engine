package chat

type RoomRepository interface {
	Add(room Room) error
	Remove(id string) error

	Get(id string) (*RoomRecord, error)
	GetAll() ([]RoomRecord, error)
}

type UserRepository interface {
	Add(user User) error
	Remove(id string) error

	Get(id string) (*UserRecord, error)
	GetAll() ([]UserRecord, error)
}
