package persist

type RoomRepository interface {
	Repository[RoomRecord, string]

	AddUser(RoomRecord, UserRecord) error
	RemoveUser(RoomRecord, UserRecord) error
}

type RoomRecord struct {
	Id    string `gorm:"primarykey"`
	Name  string
	Users []*UserRecord `gorm:"many2many:users_rooms;"` // many to many relationship
}
