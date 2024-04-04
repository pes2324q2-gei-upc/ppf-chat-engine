package chat

import (
	"bytes"
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
	Engine     ChatEngine

	// Buffered channel of outbound messages
	send chan []byte
}

// NewClient creates a new Client instance.
func NewClient(connection *websocket.Conn) *Client {
	return &Client{
		connection,
		nil,
		make(chan []byte, 256),
	}
}

func (c *Client) Send(msg []byte) {
	c.send <- msg
}

// Close closes the client connection and unregisters the client from the engine.
func (c *Client) Close() {
	c.Engine.UnregisterClient(c)
	c.Connection.Close()
	log.Println("Client disconnected")
}

// ReadPump pump messages from the websocket conn to the engine.
func (c *Client) ReadPump() {
	// defer conn closing and engine unregistering
	defer c.Close()

	// Setup client connection
	c.Connection.SetReadLimit(maxMessageSize)

	// If a message is not received within 'pongWait' duration, the read operation will return a
	// timeout error therefore asuming the client is disconnected.
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		// Set a deadline for the next read operation from the WebSocket connection.
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Loop to read messages from the WebSocket connection
	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			// Log if is unexpected close error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Send the trimmed message to the engine
		message := bytes.TrimSpace(msg)
		c.Engine.Broadcast(message)
	}
}

func (c *Client) WritePump() {
	// Start a ticker to send ping messages to the client
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	// Loop to write messages to the WebSocket connection
	for {
		select {
		case msg, ok := <-c.send:
			// Set a deadline for the next write operation to the WebSocket connection.
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Engine closed the channel
				c.Connection.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return
			}
			// Get the writer to write the message to the WebSocket connection and write the message
			writer, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(msg)

			// Handle queued messages
			for range len(c.send) {
				writer.Write(newLine)
				writer.Write(<-c.send) // Write to the WebSocket connection the msg received from the channel
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Connection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}
