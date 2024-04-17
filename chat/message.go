package chat

import (
	"encoding/json"
)

const (
	SendMessageCmd     = "SendMessage"
	GetRoomsCmd        = "GetRooms"
	GetRoomMessagesCmd = "GetRoomMessages"

	SendMessageAckContent = `{"status":"ok", "message":"sent"}`
	RoomNotFoundContent   = `{"status":"error","message":"room not found"}`
	NotImplementedContent = `{"status":"error","message":"not implemented"}`
)

type Message struct {
	MessageId string `json:"messageId"`
	Command   string `json:"command"`
	Content   string `json:"content"`
	Room      string `json:"room"`
	Sender    string `json:"sender"`
}

func (msg *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(msg),
	}

	// Unmarshal the JSON data into the auxiliary structure to avoid infinite recursion
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if msg.MessageId == "" || msg.Command == "" {
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
		// TODO why?
		if msg.Sender == "" {
			return ErrMessageMalformed
		}
	case GetRoomMessagesCmd:
		// For GetRoomMessages, all fields are required
		if msg.Content == "" || msg.Room == "" || msg.Sender == "" {
			return ErrMessageMalformed
		}
	default:
		return ErrUnknownCommand
	}
	return nil
}

func (msg *Message) Json() ([]byte, error) {
	return json.Marshal(msg)
}
