package chat

import (
	"encoding/json"
	"fmt"
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
