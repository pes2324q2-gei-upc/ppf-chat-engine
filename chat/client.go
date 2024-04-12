package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max time allowed to read the next pong msg from the peer
	pongWait   = 30 * time.Second
	pingPeriod = 10 * time.Second

	// Max time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Max message size allowed from the peer
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
)

// Client represents a WebSocket client.
type Client struct {
	Connection *websocket.Conn
	Server     *WsServer
	User       *User
	Rooms      map[string]*Room

	// Buffered channel of outbound messages
	send chan *Message
}

// Close closes the client connection and unregisters the client from the engine.
func (client *Client) Close() {}

// ReadPump pump messages from the websocket conn to the engine.
func (client *Client) ReadPump() {
	// defer conn closing and engine unregistering
	defer client.Close()

	// Setup client connection
	client.Connection.SetReadLimit(maxMessageSize)

	// If a message is not received within 'pongWait' duration, the read operation will return a
	// timeout error therefore asuming the client is disconnected.
	client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	client.Connection.SetPongHandler(func(string) error {
		// Set a deadline for the next read operation from the WebSocket connection.
		client.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Loop to read messages from the WebSocket connection
	for {
		_, msg, err := client.Connection.ReadMessage()
		if err != nil {
			// Log if is unexpected close error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Send the trimmed message to the engine
		client.HandleMessage(msg)
	}
}

func (client *Client) WritePump() {
	// Start a ticker to send ping messages to the client
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Close()
	}()
	// Loop to write messages to the WebSocket connection
	for {
		select {
		case msg, ok := <-client.send:
			// Set a deadline for the next write operation to the WebSocket connection.
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Engine closed the channel
				client.Connection.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return
			}
			// Get the writer to write the message to the WebSocket connection and write the message
			writer, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if jsonMsg, err := msg.ToJson(); err != nil {
				writer.Close()
				log.Panicln("panic: %w", err) // Panic if the message can't be marshalled to JSON
			} else {
				writer.Write(jsonMsg)
			}
			// Handle queued messages
			for range len(client.send) {
				msg = <-client.send
				if msgJson, err := msg.ToJson(); err != nil {
					writer.Close()
					log.Panicln("panic: %w", err)
				} else {
					writer.Write(newLine)
					writer.Write(msgJson)
				}
			}
			if err := writer.Close(); err != nil {
				log.Panicln("panic: %w", err)
				return
			}
		case <-ticker.C:
			if err := client.Connection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

func (client *Client) HandleMessage(msg []byte) {
	message := &Message{}
	if err := message.FromJson(msg); err != nil {
		log.Printf("error: %v", err)
		return
	}
	switch message.Command {
	case SendMessageCmd:
		client.User.Rooms[message.Room].broadcast <- message
	}

}

// NewClient creates a new Client instance.
func NewClient(connection *websocket.Conn, server *WsServer, user *User) *Client {
	return &Client{
		connection,
		server,
		user,
		make(map[string]*Room),
		make(chan *Message, 256),
	}
}
