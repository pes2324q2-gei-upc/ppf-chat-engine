package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

func TestChatEngine_RegisterClient(t *testing.T) {
	engine := chat.NewChatEngine()
	client := &chat.Client{}

	engine.RegisterClient(client)

	if len(engine.GetClients()) != 1 {
		t.Errorf("Expected 1 client, got %d", len(engine.GetClients()))
	}
}

func TestChatEngine_UnregisterClient(t *testing.T) {
	engine := chat.NewChatEngine()
	client := &chat.Client{}

	engine.RegisterClient(client)
	engine.UnregisterClient(client)

	if len(engine.GetClients()) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(engine.GetClients()))
	}
}

func TestChatEngine_Broadcast(t *testing.T) {
	// initialize 2 connections using ServeWs
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api.ServeWs(chat.DefaultChatEngine, w, r)
		}),
	)
	defer server.Close()
	con1, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatalf("(1) Failed to connect to WebSocket server: %v", err)
	}
	con2, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatalf("(2) Failed to connect to WebSocket server: %v", err)
	}

	// send a message to one of the connections
	message := []byte("test message")
	con1.WriteMessage(websocket.TextMessage, message)
	msgType, payload, error := con2.ReadMessage()
	if error != nil {
		t.Fatalf("Failed to read message from WebSocket connection: %v", error)
	}
	if msgType != websocket.TextMessage {
		t.Fatalf("Expected message type to be TextMessage, got %v", msgType)
	}
	if string(payload) != string(message) {
		t.Fatalf("Expected message to be %v, got %v", message, payload)
	}
}
