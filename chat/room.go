package chat

import db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"

type Room struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Driver *string          `json:"driver"`
	Users  map[string]*User `json:"users"`

	register   chan *User // Channel for registering/joining a user
	unregister chan *User // Channel for unregistering/leaving a user

	broadcast chan *Message // Channel for broadcasting messages to all users in the room
	close     chan bool     // Channel for closing the room
}

func (room *Room) Run() {
	for {
		select {
		case user := <-room.register:
			room.Register(user)
		case user := <-room.unregister:
			room.Unregister(user)
		case message := <-room.broadcast:
			room.Broadcast(message)
		case <-room.close:
			return
		}
	}
}

func (room *Room) Register(user *User) {
	room.Users[user.Id] = user
}

func (room *Room) Unregister(user *User) {
	delete(room.Users, user.Id)
}

func (room *Room) Broadcast(message *Message) {
	for _, user := range room.Users {
		if user.Client != nil && message.Sender != user.Id {
			user.Client.send <- message
		}
	}
}

func (room *Room) Close() {
	room.close <- true
}

func (room *Room) Empty() bool {
	return len(room.Users) == 0
}

func NewRoom(id string, name string, driver *string) *Room {
	return &Room{
		Id:         id,
		Name:       name,
		Driver:     driver,
		Users:      make(map[string]*User, 4),
		register:   make(chan *User, 2),
		unregister: make(chan *User, 2),
		broadcast:  make(chan *Message, 10),
		close:      make(chan bool),
	}
}

type RoomGateway struct {
	manager *GatewayManager
	repo    db.RoomRepository
}

func (gw *RoomGateway) Add(user *User) error {
	return nil
}

func (gw *RoomGateway) Exists(pk string) (bool, error) {
	return gw.repo.Exists(pk)
}

// Get returns a loaded User from the DB
// CAREFUL! This User rooms are nil pointers (lazy loaded)
func (gw *RoomGateway) Get(pk string) (*User, error) {
	userr, err := gw.repo.Get(pk)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:     userr.Id,
		Name:   userr.Name,
		Client: nil,
	}, nil
}
