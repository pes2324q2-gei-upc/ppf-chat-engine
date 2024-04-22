package chat

import (
	"encoding/json"
	"fmt"

	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

// User is the struct that will contain the information for a user that belongs at least to one room.
type User struct {
	Id     string           `json:"id"`
	Name   string           `json:"username"`
	Client *Client          // Web socket client that will use the user when connecting to the chat engine.
	Rooms  map[string]*Room // Map of rooms that the user belongs to.
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.Id = fmt.Sprintf("%d", aux.Id)
	u.Name = aux.Name
	return nil
}

func NewUser(id string, name string, client *Client) *User {
	return &User{
		Id:     id,
		Name:   name,
		Client: client,
		Rooms:  make(map[string]*Room),
	}
}

type UserGateway struct {
	manager *GatewayManager
	repo    db.UserRepository
}

func (gateway *UserGateway) toRecord(user *User) db.UserRecord {
	return db.UserRecord{
		Id:   user.Id,
		Name: user.Name,
	}
}

func (gateway *UserGateway) toDomain(record db.UserRecord) *User {
	// TODO implement this
	return &User{
		Id:     record.Id,
		Name:   record.Name,
		Client: nil,
		Rooms:  make(map[string]*Room),
	}
}

func (gw *UserGateway) Add(user *User) error {
	if err := gw.repo.Add(gw.toRecord(user)); err != nil {
		return err
	}
	for _, room := range user.Rooms {
		if err := gw.AddRoom(user, room); err != nil {
			return err
		}
	}
	return nil
}
func (gw *UserGateway) Exists(pk string) (bool, error) {
	return gw.repo.Exists(pk)
}

func (gw *UserGateway) Get(pk string) (*User, error) {
	record, err := gw.repo.Get(pk)
	return gw.toDomain(*record), err
}

func (gw *UserGateway) GetAll() ([]*User, error) {
	records, err := gw.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var users []*User = make([]*User, len(records))
	for _, record := range records {
		users = append(users, gw.toDomain(*record))
	}
	return users, nil
}

func (gw *UserGateway) Remove(pk string) error {
	return gw.repo.Remove(pk)
}

func (gw *UserGateway) AddRoom(user *User, room *Room) error {
	userr := gw.manager.UserGateway().toRecord(user)
	roomr := gw.manager.RoomGateway().toRecord(room)
	return gw.repo.AddRoom(userr, roomr)
}

func (gw *UserGateway) RemoveRoom(user *User, room *Room) error {
	userr := gw.manager.UserGateway().toRecord(user)
	roomr := gw.manager.RoomGateway().toRecord(room)
	return gw.repo.RemoveRoom(userr, roomr)
}
