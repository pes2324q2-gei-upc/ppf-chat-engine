package chat

type ChatEngine interface {
	RegisterClient(*Client)
	UnregisterClient(*Client)

	Broadcast([]byte)
	GetClients() []*Client
}

// WsChatEngine represents a WebSocket server.
type WsChatEngine struct {
	clients map[*Client]bool

	register   chan *Client
	unregister chan *Client

	broadcast chan []byte // broadcast channel
}

var DefaultChatEngine = NewChatEngine()

// NewChatEngine creates a new instance of WebSocketServer.
func NewChatEngine() *WsChatEngine {
	return &WsChatEngine{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run starts the WebSocket server and handles client registration and unregistration.
func (ce *WsChatEngine) Run() {
	for {
		select {
		case client := <-ce.register:
			ce.registerClient(client)
		case client := <-ce.unregister:
			ce.unregisterClient(client)
		case message := <-ce.broadcast:
			ce.broadcastClients(message)
		}
	}
}

// RegisterClient registers a client to the WebSocket server.
func (ce *WsChatEngine) registerClient(client *Client) {
	ce.clients[client] = true
	client.Engine = ce

	go client.WritePump()
	go client.ReadPump()
}

// UnregisterClient unregisters a client from the WebSocket server.
func (ce *WsChatEngine) unregisterClient(client *Client) {
	delete(ce.clients, client)
}

// RegisterClient registers a client to the WebSocket server.
func (ce *WsChatEngine) RegisterClient(client *Client) {
	ce.register <- client
}

// UnregisterClient unregisters a client from the WebSocket server.
func (ce *WsChatEngine) UnregisterClient(client *Client) {
	ce.unregister <- client
}

func (ce *WsChatEngine) GetClients() []*Client {
	clients := make([]*Client, 0, len(ce.clients))
	for client := range ce.clients {
		clients = append(clients, client)
	}
	return clients
}

// Broadcast sends a message to all clients.
func (ce *WsChatEngine) Broadcast(message []byte) {
	ce.broadcast <- message
}

// Broadcast sends a message to all clients.
func (ce *WsChatEngine) broadcastClients(message []byte) {
	for client := range ce.clients {
		client.Send(message)
	}
}
