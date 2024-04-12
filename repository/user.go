package repository

import (
	"database/sql"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

type UserRecord struct {
	Id   string
	Name string
}

type SqlUserRepository struct {
	Db *sql.DB
}

func (repo *SqlUserRepository) Add(user chat.User) error {
	_, err := repo.Db.Exec("INSERT INTO rooms (id, driver, name) VALUES (?, ?, ?)", user.Id, user.Name)
	return err
}

func (repo *SqlUserRepository) Remove(id string) error {
	_, err := repo.Db.Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}

func (repo *SqlUserRepository) Get(id string) (*UserRecord, error) {
	var user UserRecord
	row := repo.Db.QueryRow("SELECT id, name FROM user where id = ? LIMIT 1", id)

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
	}
	return &user, nil
}

func (repo *SqlUserRepository) GetAll() ([]UserRecord, error) {
	rows, err := repo.Db.Query("SELECT id, name FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]UserRecord, 0)
	for rows.Next() {
		var user UserRecord
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
