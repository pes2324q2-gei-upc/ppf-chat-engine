package chat

import (
	"database/sql"
)

type RoomRecord struct {
	Id     string
	Name   string
	Driver string
}

type SqlRoomRepository struct {
	Db *sql.DB
}

func (repo SqlRoomRepository) Add(room Room) error {
	_, err := repo.Db.Exec("INSERT INTO rooms (id, driver, name) VALUES (?, ?, ?)", room.Id, room.Driver, room.Name)
	return err
}

func (repo SqlRoomRepository) Remove(id string) error {
	_, err := repo.Db.Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}

func (repo SqlRoomRepository) Get(id string) (*RoomRecord, error) {
	var room RoomRecord
	err := repo.Db.QueryRow("SELECT id, driver, name FROM rooms WHERE id = ?", id).Scan(&room.Id, &room.Driver, &room.Name)
	return &room, err
}

func (repo SqlRoomRepository) GetAll() ([]RoomRecord, error) {
	rows, err := repo.Db.Query("SELECT id, driver, name FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rooms := make([]RoomRecord, 0, 10)
	for rows.Next() {
		var room RoomRecord
		if err := rows.Scan(&room.Id, &room.Driver, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
