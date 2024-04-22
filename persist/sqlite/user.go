package sqlite

import (
	"database/sql"
	"errors"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

var ErrInvalidRecordType = errors.New("invalid record type")

func parseUserRows(rows *sql.Rows) ([]*persist.UserRecord, error) {
	var users []*persist.UserRecord
	for rows.Next() {
		user := &persist.UserRecord{}
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

type SqlUserRepository struct {
	Db *sql.DB
}

func (repo SqlUserRepository) Exists(id string) (bool, error) {
	rows, err := repo.Db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo SqlUserRepository) Add(user persist.UserRecord) error {
	_, err := repo.Db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.Id, user.Name)
	return err
}

func (repo SqlUserRepository) Remove(id string) error {
	_, err := repo.Db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (repo SqlUserRepository) Get(id string) (*persist.UserRecord, error) {
	row := repo.Db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
	user := &persist.UserRecord{}
	err := row.Scan(user.Id, user.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo SqlUserRepository) GetAll() ([]*persist.UserRecord, error) {
	rows, err := repo.Db.Query("SELECT id, name FROM users", nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseUserRows(rows)
}

func (repo SqlUserRepository) AddRoom(user persist.UserRecord, room persist.RoomRecord) error {
	_, err := repo.Db.Exec("INSERT INTO room_user (room, user) VALUES (?, ?)", room.Id, user.Id)
	return err
}

func (repo SqlUserRepository) RemoveRoom(user persist.UserRecord, room persist.RoomRecord) error {
	_, err := repo.Db.Exec("DELETE FROM room_user WHERE room = ? AND user = ?", room.Id, user.Id)
	return err
}
