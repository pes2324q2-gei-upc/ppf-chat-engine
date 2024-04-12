package chat

type User struct {
	Id     string
	Name   string
	Client *Client
	Rooms  map[string]*Room
}

func NewUser(id, name string, client *Client) *User {
	return &User{
		Id:     id,
		Name:   name,
		Client: client,
		Rooms:  make(map[string]*Room),
	}
}
