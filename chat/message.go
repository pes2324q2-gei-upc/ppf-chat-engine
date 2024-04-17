package chat

import (
	"encoding/json"
	"fmt"
)

const (
	SendMessageCmd = "SendMessage"
	GetRoomsCmd    = "GetRooms"
	GetMessagesCmd = "GetMessages"
)

type Message struct {
	Command string `json:"command"`
	Content string `json:"content"`
	Room    string `json:"room"`
	Sender  string `json:"sender"`
}

func (msg *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(msg),
	}

	// Unmarshal the JSON data into the auxiliary structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if msg.Command == "" {
		return ErrMessageMalformed
	}
	// Check if non-optional fields are empty based on the command
	switch msg.Command {
	case SendMessageCmd:
		// For SendMessage, all fields are required
		if msg.Content == "" || msg.Room == "" || msg.Sender == "" {
			return ErrMessageMalformed
		}
	case GetRoomsCmd:
		// For GetRooms, only sender is required
		if msg.Sender == "" {
			return ErrMessageMalformed
		}
	case GetMessagesCmd:
		// For GetMessages, all fields are required
		if msg.Content == "" || msg.Room == "" || msg.Sender == "" {
			return ErrMessageMalformed
		}
	default:
		return fmt.Errorf("unknown command: %s", msg.Command)
	}
	return nil
}

func (msg *Message) Json() ([]byte, error) {
	return json.Marshal(msg)
}
