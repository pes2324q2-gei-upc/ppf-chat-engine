package sqlite

import (
	"database/sql"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

func parseRows(rows *sql.Rows) ([]*persist.MessageRecord, error) {
	var messages []*persist.MessageRecord
	for rows.Next() {
		msg := &persist.MessageRecord{}
		err := rows.Scan(&msg.Room, &msg.Sender, &msg.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

type SqlMessageRepository struct {
	Db *sql.DB
}

func (repo *SqlMessageRepository) Exists(id persist.MessageKey) (bool, error) {
	rows, err := repo.Db.Query("SELECT * FROM messages WHERE room = ? AND sender = ?", id.Room, id.Sender)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo *SqlMessageRepository) Add(msg persist.MessageRecord) error {
	_, err := repo.Db.Exec("INSERT INTO messages (room, sender, content) VALUES (?, ?, ?, ?)", msg.Room, msg.Sender, msg.Content)
	return err
}

func (repo *SqlMessageRepository) Remove(pk persist.MessageKey) error {
	_, err := repo.Db.Exec("DELETE FROM messages WHERE room = ? AND sender = ?", pk.Room, pk.Sender)
	return err
}

func (repo *SqlMessageRepository) Get(id persist.MessageKey) (*persist.MessageRecord, error) {
	row := repo.Db.QueryRow("SELECT * FROM messages WHERE room = ? AND sender = ?", id.Room, id.Sender)
	msg := persist.MessageRecord{}
	err := row.Scan(&msg.Room, &msg.Sender, &msg.Content)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (repo *SqlMessageRepository) GetAll() ([]*persist.MessageRecord, error) {
	rows, err := repo.Db.Query("SELECT * FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseRows(rows)
}

func (repo *SqlMessageRepository) GetByRoom(room string) ([]*persist.MessageRecord, error) {
	rows, err := repo.Db.Query("SELECT * FROM messages WHERE room = ?", room)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseRows(rows)
}

func (repo *SqlMessageRepository) GetBySender(sender string) ([]*persist.MessageRecord, error) {
	rows, err := repo.Db.Query("SELECT * FROM messages WHERE sender = ?", sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseRows(rows)
}
