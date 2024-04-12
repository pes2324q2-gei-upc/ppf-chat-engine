package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatEngine represents the engine that manages the chat rooms and users.
type ChatEngine struct {
	HttpClient *http.Client
	Server     *WsServer        // Server represents the WebSocket server.
	Users      map[string]*User // Users represents the map of users in the chat engine.
	Rooms      map[string]*Room // Rooms represents the map of rooms in the chat engine.
	UserRepo   UserRepository   // UserRepo represents the repository for users.
	RoomRepo   RoomRepository   // RoomRepo represents the repository for rooms.
}

// CloseRoom closes the specified room and removes it from the chat engine.
func (engine *ChatEngine) CloseRoom(id string) {
	room := engine.Rooms[id]
	room.close <- true

	// If a user is not in any room, close the connection and delete it.
	for _, user := range room.Users {
		if len(user.Rooms) == 0 {
			engine.Server.unregister <- user.Client
			delete(engine.Users, user.Id)
		}
	}
	delete(engine.Rooms, id)
}

// ConnectUser connects a user to the chat engine with the specified ID.
// The user must be loaded in the chat engine before connecting.
func (engine *ChatEngine) ConnectUser(id string, w http.ResponseWriter, r *http.Request) error {
	client := engine.Server.OpenConnection(w, r)
	if user, ok := engine.Users[id]; !ok {
		return ErrUserNotLoaded
	} else {
		client.User = user
	}
	// Register the client to the server.
	engine.Server.register <- client
	return nil
}

func (engine *ChatEngine) Exists(id string) bool {
	_, ok := engine.Users[id]
	return ok
}

// JoinRoom joins a user to the specified room in the chat engine.
// The user must be loaded in the chat engine before connecting.
func (engine *ChatEngine) JoinRoom(room string, userId string) error {
	// If the user is not in the engine, create it.
	if user, ok := engine.Users[userId]; !ok {
		return ErrUserNotLoaded
	} else {
		engine.Rooms[room].register <- user
	}
	return nil
}

// LeaveRoom removes a user from the specified room in the chat engine.
// If the room ends up empty, it will be closed.
// If the user ends up in no rooms, the connection will be closed and the user will be deleted.
func (engine *ChatEngine) LeaveRoom(roomId string, userId string) error {
	user, ok := engine.Users[userId]
	if !ok {
		return ErrUserNotLoaded
	}
	room, ok := engine.Rooms[roomId]
	if !ok {
		return ErrRoomNotFound
	}
	room.unregister <- user
	if room.Empty() {
		engine.CloseRoom(roomId)
	}
	// If the user is in no rooms, close the connection and delete it.
	if len(user.Rooms) == 0 {
		engine.Server.unregister <- user.Client
		delete(engine.Users, userId)
	}
	return nil
}

func (engine *ChatEngine) LoadUser(id string) error {
	// Make a get request to the UserAPI to retrieve the user
	r, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://user-api/drivers/%s", id),
		nil,
	)
	response, err := engine.HttpClient.Do(r) // TODO handle error
	if err != nil || response.StatusCode != http.StatusOK {
		return ErrUserApiRequestFailed
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body) // TODO handle error
	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return ErrUserUnmarshalFailed
	}
	engine.Users[id] = &user
	return nil
}

// OpenRoom creates a new room with the specified ID, name and driver user, and adds it to the engine.
// The driver user must be loaded in the chat engine before opening the room.
func (engine *ChatEngine) OpenRoom(id string, name string, driverId string) error {
	// If the driver is not in the engine, return an error.
	driver, ok := engine.Users[driverId]
	if !ok {
		return ErrUserNotLoaded
	}
	room := NewRoom(id, name, &driver.Id)
	engine.Rooms[id] = room

	go room.Run()
	return nil
}

// NewChatEngine creates a new instance of ChatEngine.
func NewChatEngine(conn *sql.DB) *ChatEngine {
	engine := &ChatEngine{
		HttpClient: &http.Client{},
		UserRepo:   SqlUserRepository{Db: conn},
		RoomRepo:   SqlRoomRepository{Db: conn},
	}
	engine.Server = NewWsServer(engine)

	return engine
}
