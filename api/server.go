package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

// RegisterHandlers registers the HTTP request handlers.
func RegisterHandlers() {
	log.Println("Registering handlers")
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/ws", WsHandler)
}

// RootHandler handles the root HTTP request, always sends 200 status code and an empty body.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// WsHandler handles the WebSocket request.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ServeWs(chat.DefaultChatEngine, w, r)
}

// ServeWs upgrades the HTTP connection to a WebSocket connection and creates a new server client for the connection.
func ServeWs(engine chat.ChatEngine, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := chat.NewClient(conn)
	engine.RegisterClient(client)
	log.Println("New client connected")
}
