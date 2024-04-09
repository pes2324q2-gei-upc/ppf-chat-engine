package chat

type UserId string

type User struct {
	Id     UserId
	Name   string
	Client *Client
	Rooms  map[RoomId]*Room
}
