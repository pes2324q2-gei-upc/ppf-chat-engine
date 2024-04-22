package persist

import "time"

type MessageKey struct {
	Room   string
	Sender string
}

type MessageRepository interface {
	Repository[MessageRecord, MessageKey]

	GetByRoom(room string) ([]*MessageRecord, error)
	GetBySender(sender string) ([]*MessageRecord, error)
}

type MessageRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	Room    RoomRecord // Belongs To relation
	Sender  UserRecord // Belongs To relation
	Content string
}
