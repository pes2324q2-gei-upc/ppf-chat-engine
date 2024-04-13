package repository

import "github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"

type RoomRepository interface {
	Add(room chat.Room) error
	Remove(id string) error

	Get(id string) (RoomRecord, error)
	GetAll() ([]RoomRecord, error)
}

type UserRepository interface {
	Add(user chat.User) error
	Remove(id string) error

	Get(id string) (UserRecord, error)
	GetAll() ([]UserRecord, error)
}
