package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	Engine     *ChatEngine
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			log.Printf("info: registering client %s", client.User.Id)
			server.Clients[client] = true
		case client := <-server.unregister:
			if _, ok := server.Clients[client]; ok {
				delete(server.Clients, client)
				close(client.send)
			}
		}
	}
}

func (server *WsServer) OpenConnection(w http.ResponseWriter, r *http.Request) *Client {
	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	client := NewClient(conn, server, nil)
	go client.ReadPump()
	go client.WritePump()
	return client
}

func NewWsServer(engine *ChatEngine) *WsServer {
	return &WsServer{
		Engine:     engine,
		Clients:    make(map[*Client]bool, 64),
		register:   make(chan *Client, 2),
		unregister: make(chan *Client, 2),
	}
}
