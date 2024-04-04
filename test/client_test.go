package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	api "github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

func TestNewClient(t *testing.T) {
	conn := &websocket.Conn{}
	engine := chat.DefaultChatEngine

	client := chat.NewClient(conn)
	engine.RegisterClient(client)

	if client.Connection != conn {
		t.Errorf("NewClient() failed: expected connection to be %v, got %v", conn, client.Connection)
	}

	if client.Engine != engine {
		t.Errorf("NewClient() failed: expected engine to be %v, got %v", engine, client.Engine)
	}
}

// func TestServeWs(t *testing.T) {
// 	engine := &chat.WsChatEngine{}
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		chat.ServeWs(engine, w, r)
// 	})

// 	// Create a test server
// 	server := httptest.NewServer(handler)
// 	defer server.Close()

// 	// Create a WebSocket connection to the test server
// 	wsURL := "ws" + server.URL[4:]
// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		t.Fatalf("Failed to establish WebSocket connection: %v", err)
// 	}
// 	defer conn.Close()

// 	// Verify that a new client is registered with the engine
// 	clients := engine.GetClients()
// 	if len(clients) != 1 {
// 		t.Errorf("ServeWs() failed: expected 1 client to be registered, got %d", len(clients))
// 	}
// }

// func TestClient_Close(t *testing.T) {
// 	conn := &websocket.Conn{}
// 	engine := &chat.WsChatEngine{}
// 	client := chat.NewClient(conn)
// 	engine.RegisterClient(client)

// 	// Call the Close method
// 	client.Close()

// 	// Verify that the client is unregistered from the engine
// 	clients := engine.GetClients()
// 	if len(clients) != 0 {
// 		t.Errorf("Close() failed: expected 0 clients to be registered, got %d", len(clients))
// 	}
// }

func TestClient_ReadPump(t *testing.T) {
	// ReadPump pumps messages from the WebSocket connection to the client's receive channel.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(chat.DefaultChatEngine, w, r)
	}))
	defer s.Close()
	wsConn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(s.URL, "http"), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer wsConn.Close()

	client := chat.DefaultChatEngine.GetClients()[0]
	// Create a test message
	message := []byte("test message")
	// Send the test message to the client
	if err := wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatalf("Failed to send test message: %v", err)
	}
	// Call the ReadPump method
	go client.ReadPump()
	// Wait for the message to be received by the client
	time.Sleep(time.Millisecond)
}

func TestClient_WritePump(t *testing.T) {
	// WritePump pumps messages from the client's send channel to the WebSocket connection.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(chat.DefaultChatEngine, w, r)
	}))
	defer s.Close()
	wsConn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(s.URL, "http"), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer wsConn.Close()

	client := chat.DefaultChatEngine.GetClients()[0]
	// Create a test message
	message := []byte("test message")
	// Create a channel to receive messages sent by the client
	received := make(chan []byte)

	// Start a goroutine to handle messages received by the client
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				break
			}
			received <- message
		}
	}()
	// Call the WritePump method
	go client.WritePump()
	// Send the test message to the client
	client.Send(message)
	// Wait for the message to be sent by the client
	time.Sleep(time.Millisecond)

	// Verify that the message was sent by the client
	sent := <-received
	if string(sent) != string(message) {
		t.Errorf("WritePump() failed: expected message to be sent, got %s", string(sent))
	}
}
