package chat

import (
	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

type MessageKey struct {
	room   string
	sender string
}

type Message struct {
	Content string `json:"content"`
	Room    Room   `json:"room"`
	Sender  User   `json:"sender"`
}

// MessageGateway acts as a data mapper between the domain layer and de data layer, transforming the data from the database into domain objects and vice versa.
type MessageGateway struct {
	repo db.Repository[db.MessageKey, db.Message]
}

func (gw MessageGateway) MessageRecordToMessage(record db.Message) Message {
	r := RoomGateway{}.RoomRecordToRoom(record.Room)
	u := UserGateway{}.UserRecordToUser(record.Sender)
	return Message{
		Content: record.Content,
		Room:    r,
		Sender:  u,
	}
}

func (gw MessageGateway) MessageToMessageRecord(msg Message) db.Message {
	r := RoomGateway{}.RoomToRoomRecord(msg.Room)
	u := UserGateway{}.UserToUserRecord(msg.Sender)
	return db.Message{
		Room:    r,
		Sender:  u,
		Content: msg.Content,
	}
}

func (gw MessageGateway) Exists(pk MessageKey) bool {
	key := db.MakeMessageKey(pk.room, pk.sender)
	return gw.repo.Exists(key)
}

func (gw MessageGateway) Create(room *Message) error {
	msgr := gw.MessageToMessageRecord(*room)
	return gw.repo.Create(msgr)
}

func (gw MessageGateway) Read(pk MessageKey) *Message {
	key := db.MakeMessageKey(pk.room, pk.sender)
	msgr := gw.repo.Read(key)
	room := gw.MessageRecordToMessage(msgr)
	return &room

}

func (gw MessageGateway) ReadAll() []*Message {
	msgrs := gw.repo.ReadAll()
	rooms := make([]*Message, 0)
	for _, u := range msgrs {
		room := gw.MessageRecordToMessage(u)
		rooms = append(rooms, &room)
	}
	return rooms
}

func (gw MessageGateway) Update(pk MessageKey, user *Message) error {
	return nil
}

func (gw MessageGateway) Delete(pk MessageKey) {
	key := db.MakeMessageKey(pk.room, pk.sender)
	gw.repo.Delete(key)
}
