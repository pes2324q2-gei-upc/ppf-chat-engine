package chat

import (
	"encoding/json"
	"fmt"

	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

type User struct {
	Id     string `json:"id"`
	Name   string `json:"username"`
	Client *Client
	Rooms  map[string]*Room
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

type UserGateway struct{}

func (gateway *UserGateway) ToRecord(user User) db.Record {
	return &db.UserRecord{
		Id: user.Id,
	}
}

func (gateway *UserGateway) ToDomain(record db.Record) User {
	return User{
		Id:   record.Pk().(string),
		Name: record.(*db.UserRecord).Name,
	}
}
