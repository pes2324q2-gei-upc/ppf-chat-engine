package sqlite

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
	CreateRoomUserTable = `CREATE TABLE IF NOT EXISTS room_user (
		room VARCHAR(128) NOT NULL,
		user VARCHAR(128) NOT NULL
		FOREIGN KEY (room) REFERENCES rooms(id)
		FOREIGN KEY (user) REFERENCES users(id)
	)`
	CreateMessagesTable = `CREATE TABLE IF NOT EXISTS messages (
		ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		room VARCHAR(128) NOT NULL,
		sender VARCHAR(128) NOT NULL,
		content TEXT NOT NULL
		PRIMARY KEY (ts, room, sender)
		FOREING KEY (room) REFERENCES rooms(id)
		FOREING KEY (sender) REFERENCES users(id)
	);`
)

// InitDB initializes the database with the specified driver and source.
func Init(driver string, source string) *sql.DB {
	log.Printf("info: init %s database at %s", driver, source)

	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateRoomsTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateUsersTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateRoomUserTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(CreateMessagesTable); err != nil {
		log.Fatal(err)
	}
	return db
}
