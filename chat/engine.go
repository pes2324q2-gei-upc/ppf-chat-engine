package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/auth"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/config"
	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/persist/sqlite"
)

var (
	RepoUserGw = &UserGateway{}
	// RepoRoomGw = &RoomGateway{}
	// RepoMsgGw  = &MessageGateway{}
)

// ChatEngine represents the engine that manages the chat rooms and users.
type ChatEngine struct {
	Configuration *config.Configuration // Configuration represents the configuration of the chat engine.
	HttpClient    *http.Client
	Server        *WsServer            // Server represents the WebSocket server.
	Users         map[string]*User     // Users represents the map of users in the chat engine.
	Rooms         map[string]*Room     // Rooms represents the map of rooms in the chat engine.
	UserRepo      db.UserRepository    // UserRepo represents the repository for users.
	RoomRepo      db.RoomRepository    // RoomRepo represents the repository for rooms.
	MessageRepo   db.MessageRepository // MessageRepo represents the repository for messages.
}

// CloseRoom closes the specified room and removes it from the chat engine.
func (engine *ChatEngine) CloseRoom(id string) error {
	room, ok := engine.Rooms[id]
	if !ok {
		return ErrRoomNotFound
	}
	room.close <- true
	// If a user is not in any room, close the connection and delete it.
	for _, user := range room.Users {
		if len(user.Rooms) == 0 {
			engine.Server.unregister <- user.Client
			delete(engine.Users, user.Id)
		}
	}
	delete(engine.Rooms, id)
	return nil
}

// ConnectUser connects a user to the chat engine with the specified ID.
// The user must be loaded in the chat engine before connecting.
func (engine *ChatEngine) ConnectUser(id string, w http.ResponseWriter, r *http.Request) error {
	log.Printf("info: connecting user %s", id)
	user, ok := engine.Users[id]
	if !ok {
		log.Printf("error: user %s not found", id)
		return ErrUserNotFound
	}
	client := engine.Server.OpenConnection(w, r)
	client.User = user
	user.Client = client
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
	log.Printf("info: joining user %s to room %s", userId, room)
	if user, ok := engine.Users[userId]; !ok {
		log.Printf("error: user %s not found", userId)
		return ErrUserNotFound
	} else if r, ok := engine.Rooms[room]; !ok {
		log.Printf("error: room %s not found", room)
		return ErrRoomNotFound
	} else {
		user.Rooms[room] = r
		r.register <- user
	}
	return nil
}

// LeaveRoom removes a user from the specified room in the chat engine.
// If the room ends up empty, it will be closed.
// If the user ends up in no rooms, the connection will be closed and the user will be deleted.
func (engine *ChatEngine) LeaveRoom(roomId string, userId string) error {
	user, ok := engine.Users[userId]
	if !ok {
		return ErrUserNotFound
	}
	room, ok := engine.Rooms[roomId]
	if !ok {
		return ErrRoomNotFound
	}
	room.unregister <- user
	// If the user is in no rooms, close the connection and delete it.
	if len(user.Rooms) == 0 {
		engine.Server.unregister <- user.Client
		delete(engine.Users, userId)
	}
	return nil
}

func (engine *ChatEngine) LoadUser(id string) error { // IMPROVE error handling
	log.Printf("info: loading user %s", id)

	usrUrl := engine.Configuration.UserApiUrl.JoinPath("drivers", id)
	r, _ := http.NewRequest(
		http.MethodGet,
		usrUrl.String(),
		nil,
	)
	r.Header.Add("Authorization", fmt.Sprintf("Token %s", engine.Configuration.Credentials.Token()))
	response, err := engine.HttpClient.Do(r)

	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("error: could not load user %s: %v", id, ErrUserApiRequestFailed)
		return ErrUserApiRequestFailed
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	user := NewUser("", "", nil)
	if err = json.Unmarshal(body, user); err != nil {
		log.Printf("error: could not load user %s: %v", id, ErrUserUnmarshalFailed)
		return ErrUserUnmarshalFailed
	}
	engine.Users[id] = user
	return nil
}

// OpenRoom creates a new room with the specified ID, name and driver user, and adds it to the engine.
// The driver user must be loaded in the chat engine before opening the room.
func (engine *ChatEngine) OpenRoom(id string, name string, driver string) error {
	log.Printf("info: opening room %s", id)
	user, ok := engine.Users[driver]
	if !ok {
		log.Printf("error: driver %s not found", driver)
		return ErrUserNotFound
	}
	room := NewRoom(id, name, &user.Id)
	engine.Rooms[id] = room

	go room.Run()
	engine.JoinRoom(id, driver)
	return nil
}

func NewDefaultChatEngine(db *sql.DB) (*ChatEngine, error) {
	useUrl, _ := url.Parse(config.GetEnv("USER_API_URL", "http://localhost:8081"))
	routeUrl, _ := url.Parse(config.GetEnv("ROUTE_API_URL", "http://localhost:8080"))

	credentials := &auth.UserApiCredentials{
		AuthUrl:  useUrl,
		Email:    config.GetEnv("PPF_MAIL", "admin@ppf.com"),
		Password: config.GetEnv("PPF_PASS", "chatengine"),
	}
	if err := credentials.Login(); err != nil {
		return nil, err
	}
	configuration := &config.Configuration{
		Debug:       config.GetEnv("DEBUG", "false") == "true",
		UserApiUrl:  useUrl,
		RouteApiUrl: routeUrl,
		Credentials: credentials,
	}
	engine := &ChatEngine{
		Configuration: configuration,
		HttpClient:    http.DefaultClient,
		Users:         make(map[string]*User),
		Rooms:         make(map[string]*Room),
		UserRepo:      &sqlite.SqlUserRepository{Db: db},
		RoomRepo:      &sqlite.SqlRoomRepository{Db: db},
		MessageRepo:   &sqlite.SqlMessageRepository{Db: db},
	}
	return engine, nil
}
