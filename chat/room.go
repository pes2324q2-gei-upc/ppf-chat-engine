package chat

type Room struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Driver     *string          `json:"driver"`
	Users      map[string]*User `json:"users"`
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
		Users:      make(map[string]*User, 1),
		register:   make(chan *User, 2),
		unregister: make(chan *User, 2),
		broadcast:  make(chan *Message, 2),
		connect:    make(chan *Client, 2),
		close:      make(chan bool),
	}
}
