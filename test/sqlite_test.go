package chat

import (
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

func TestInitDB(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db := chat.InitDB("sqlite3", "file::memory:")

	// Check if the database connection is valid
	t.Run("Check database connection", func(t *testing.T) {
		if err := db.Ping(); err != nil {
			t.Errorf("Expected database connection to be valid, got error: %v", err)
		}
	})

	// Check if the rooms table exists
	t.Run("Check rooms table", func(t *testing.T) {
		_, err := db.Exec("SELECT * FROM rooms")
		if err != nil {
			t.Errorf("Expected rooms table to exist, got error: %v", err)
		}
	})

	// Check if the users table exists
	t.Run("Check users table", func(t *testing.T) {
		_, err := db.Exec("SELECT * FROM users")
		if err != nil {
			t.Errorf("Expected users table to exist, got error: %v", err)
		}
	})

	// Close the database connection
	t.Run("Close database connection", func(t *testing.T) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	})
}
