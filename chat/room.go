package chat

type RoomId string

type Room struct {
	Id         RoomId
	Name       string
	Driver     *UserId
	Users      map[UserId]*User
	register   chan *User
	unregister chan *User
	broadcast  chan *Message
	connect    chan *Client
	close      chan bool
}

func (room *Room) Run() {
	for {
		select {
		case client := <-room.connect:
			room.Users[client.User.Id].Client = client
		case user := <-room.register:
			room.Users[user.Id] = user
		case user := <-room.unregister:
			delete(room.Users, user.Id)
		case message := <-room.broadcast:
			for _, user := range room.Users {
				user.Client.send <- message
			}
		case <-room.close:
			return
		}
	}
}

func (room *Room) Empty() bool {
	return len(room.Users) == 0
}

func NewRoom(id RoomId, name string, driver *UserId) *Room {
	return &Room{
		Id:       id,
		Name:     name,
		Driver:   driver,
		Users:    make(map[UserId]*User, 1),
		register: make(chan *User, 2),
	}
}
