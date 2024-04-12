package chat

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CreateRoomsTable = `CREATE TABLE IF NOT EXISTS rooms (
		id VARCHAR(128) PRIMARY KEY,
		driver VARCHAR(128) NOT NULL,
		name VARCHAR(128) NOT NULL
	);`
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(128) PRIMARY KEY,
		name VARCHAR(128) NOT NULL
	);`
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateRoomsTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateUsersTable); err != nil {
		log.Fatal(err)
	}
	return db
}
