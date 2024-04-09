package chat

import "encoding/json"

const (
	SendMessageCmd = "send_message"
)

type Message struct {
	Command string `json:"command"`
	Content string `json:"content"`
	Room    RoomId `json:"room"`
	Sender  UserId `json:"sender"`
}

func (m *Message) FromJson(jsonData []byte) error {
	return json.Unmarshal(jsonData, m)
}

func (m *Message) ToJson() ([]byte, error) {
	return json.Marshal(m)
}
