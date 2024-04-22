package sqlite

import (
	"database/sql"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

func parseRoomRows(rows *sql.Rows) ([]*persist.RoomRecord, error) {
	var rooms []*persist.RoomRecord
	for rows.Next() {
		room := &persist.RoomRecord{}
		err := rows.Scan(&room.Id, &room.Name)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

type SqlRoomRepository struct {
	Db *sql.DB
}

func (repo SqlRoomRepository) Exists(id string) (bool, error) {
	rows, err := repo.Db.Query("SELECT * FROM rooms WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo SqlRoomRepository) Add(room persist.RoomRecord) error {
	_, err := repo.Db.Exec("INSERT INTO rooms (id, name) VALUES (?, ?)", room.Id, room.Name)
	return err
}

func (repo SqlRoomRepository) Remove(id string) error {
	_, err := repo.Db.Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}

func (repo SqlRoomRepository) Get(id string) (*persist.RoomRecord, error) {
	row := repo.Db.QueryRow("SELECT id, name FROM rooms WHERE id = ?", id)
	room := &persist.RoomRecord{}
	err := row.Scan(room.Id, room.Name)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (repo SqlRoomRepository) GetAll() ([]*persist.RoomRecord, error) {
	rows, err := repo.Db.Query("SELECT id, name FROM rooms", nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseRoomRows(rows)
}

func (repo SqlRoomRepository) AddUser(room persist.RoomRecord, user persist.UserRecord) error {
	_, err := repo.Db.Exec("INSERT INTO room_user (room, user) VALUES (?, ?)", room.Id, user.Id)
	return err
}

func (repo SqlRoomRepository) RemoveUser(room persist.RoomRecord, user persist.UserRecord) error {
	_, err := repo.Db.Exec("DELETE FROM room_user WHERE room = ? AND user = ?", user.Id, room.Id)
	return err
}
