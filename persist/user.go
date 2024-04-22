package persist

type UserRepository interface {
	Repository[UserRecord, string]

	AddRoom(UserRecord, RoomRecord) error
	RemoveRoom(UserRecord, RoomRecord) error
}

type UserRecord struct {
	Id    string `gorm:"primarykey"`
	Name  string
	Rooms []*RoomRecord `gorm:"many2many:users_rooms;"` // many to many relationship
}
