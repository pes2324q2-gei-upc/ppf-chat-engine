package persist

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
	Room    string
	Sender  string
	Content string
}
