package chat

import (
	"encoding/json"
	"fmt"

	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

// User is the struct that will contain the information for a user that belongs at least to one room.
type User struct {
	Id     string  `json:"id"`
	Name   string  `json:"username"`
	Client *Client // Web socket client that will use the user when connecting to the chat engine.
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
	}
}

type UserGateway struct {
	manager *GatewayManager
	repo    db.UserRepository
}

func (gw *UserGateway) UserRecordToUser(record db.UserRecord) User {
	return User{
		Id:     record.Pk(),
		Name:   record.GetName(),
		Client: nil,
	}
}

func (gw *UserGateway) UserToUserRecord(user User) db.UserRecord {
	return db.UserRecord{
		Id:   user.Id,
		Name: user.Name,
	}
}

func (gw *UserGateway) Add(user *User) error {
	if err := gw.repo.Add(gw.UserToUserRecord(*user)); err != nil {
		return err
	}
	return nil
}

func (gw *UserGateway) Exists(pk string) (bool, error) {
	return gw.repo.Exists(pk)
}

// Get returns a loaded User from the DB
// CAREFUL! This User rooms are nil pointers (lazy loaded)
func (gw *UserGateway) Get(pk string) (*User, error) {
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

func (gw *UserGateway) GetAll() ([]*User, error) {
	return nil, nil
}

func (gw *UserGateway) Remove(pk string) error {
	return gw.repo.Remove(pk)
}
