package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/auth"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/config"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
	"gorm.io/gorm"
)

// ChatEngine represents the engine that manages the chat rooms and users.
type ChatEngine struct {
	GatewayManager
	config.Configuration // Configuration represents the configuration of the chat engine.
	HttpClient           *http.Client
	Server               *WsServer        // Server represents the WebSocket server.
	Users                map[string]*User // Users represents the map of users in the chat engine.
	Rooms                map[string]*Room // Rooms represents the map of rooms in the chat engine.
}

func (engine *ChatEngine) AddUser(user *User) {
	engine.Users[user.Id] = user
}

// CloseRoom closes the specified room and removes it from the chat engine.
func (engine *ChatEngine) CloseRoom(id string) error {
	room, ok := engine.Rooms[id]
	if !ok {
		return ErrRoomNotFound
	}
	room.close <- true
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

// InitUser initializes a user in the chat engine by:
// 1. Requesting the user data at the UserAPI
// 2. Requesting the routes that belong to the user (as passenger or driver) at the RouteAPI
// 3. Adding the user to the engine
// 4. Opening a room for each route
// 5. Requesting all users from a route
// 6. Adding those users to the engine

func (engine *ChatEngine) InitUser(id string) error {
	// Get the user from UserAPI
	user, err := engine.RequestUser(id)
	if err != nil {
		return err
	}
	// Get the routes from a user
	routes, err := engine.RequestUserRoutes(id)
	if err != nil {
		return err
	}
	engine.AddUser(user)
	// for each route open a room
	for _, route := range routes {
		// but request driver user first and join it
		driver, err := engine.RequestUser(route.Driver.Id)
		if err != nil {
			return err
		}
		engine.AddUser(driver)
		engine.OpenRoom(route.Id, route.Name, driver.Id)
		engine.JoinRoom(driver.Id, route.Id)
		// then request all passengers from a route
		users, err := engine.RequestRoutePassengers(route.Id)
		if err != nil {
			return err
		}
		for _, user := range users {
			// and join them to the room
			engine.AddUser(user)
			engine.JoinRoom(route.Id, user.Id)
		}
	}
	return nil
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

func (engine *ChatEngine) RequestRoutePassengers(id string) ([]*User, error) {
	return nil, nil
}

// RequestUser loads the user by getting it from the DB and, if it does not exist, from the user API.
func (engine *ChatEngine) RequestUser(id string) (*User, error) {
	log.Printf("info: loading user %s", id)
	usrUrl := engine.Configuration.UserApiUrl.JoinPath("drivers", id)
	userReq, _ := http.NewRequest(
		http.MethodGet,
		usrUrl.String(),
		nil,
	)
	userReq.Header.Add("Authorization", fmt.Sprintf("Token %s", engine.Configuration.Credentials.Token()))
	response, err := engine.HttpClient.Do(userReq)

	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("error: could not load user %s: %v", id, ErrUserApiRequestFailed)
		return nil, ErrUserApiRequestFailed
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	user := NewUser("", "", nil)
	if err = json.Unmarshal(body, user); err != nil {
		log.Printf("error: could not load user %s: %v", id, ErrUserUnmarshalFailed)
		return nil, ErrUserUnmarshalFailed
	}
	return user, nil
}

func (engine *ChatEngine) RequestUserRoutes(id string) ([]*Route, error) {
	return nil, nil
}

// NewChatEngine creates a new chat engine with the intended application defaults.
func NewDefaultChatEngine(db *gorm.DB) (*ChatEngine, error) {
	useUrl, _ := url.Parse(config.GetEnv("USER_API_URL", "http://localhost:8081"))
	routeUrl, _ := url.Parse(config.GetEnv("ROUTE_API_URL", "http://localhost:8080"))

	credentials := auth.UserApiCredentials{
		AuthUrl:  useUrl,
		Email:    config.GetEnv("PPF_MAIL", "admin@ppf.com"),
		Password: config.GetEnv("PPF_PASS", "chatengine"),
	}
	if err := credentials.Login(); err != nil {
		return nil, err
	}
	configuration := config.Configuration{
		Debug:       config.GetEnv("DEBUG", "false") == "true",
		UserApiUrl:  *useUrl,
		RouteApiUrl: *routeUrl,
		Credentials: credentials,
	}
	gwm := GatewayManager{
		userGw: UserGateway{
			repo: persist.UserRepository{DB: db},
		},
		roomGw: RoomGateway{
			repo: persist.RoomRepository{DB: db},
		},
		msgGw: MessageGateway{
			repo: persist.MessageRepository{DB: db},
		},
	}
	gwm.userGw.manager = &gwm
	gwm.roomGw.manager = &gwm
	gwm.msgGw.manager = &gwm

	engine := &ChatEngine{
		Configuration:  configuration,
		HttpClient:     http.DefaultClient,
		GatewayManager: gwm,
		Users:          make(map[string]*User),
		Rooms:          make(map[string]*Room),
	}
	return engine, nil
}
