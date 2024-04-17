package chat

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ChatEngine represents the engine that manages the chat rooms and users.
type ChatEngine struct {
	Configuration *Configuration // Configuration represents the configuration of the chat engine.
	HttpClient    *http.Client
	Server        *WsServer                    // Server represents the WebSocket server.
	Users         map[string]*User             // Users represents the map of users in the chat engine.
	Rooms         map[string]*Room             // Rooms represents the map of rooms in the chat engine.
	UserRepo      Repository[User, UserRecord] // UserRepo represents the repository for users.
	RoomRepo      Repository[Room, RoomRecord] // RoomRepo represents the repository for rooms.
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
func (engine *ChatEngine) JoinRoom(room string, user string) error {
	// If the user is not in the engine, create it.
	log.Printf("info: joining user %s to room %s", user, room)
	if user, ok := engine.Users[user]; !ok {
		log.Printf("error: user %s not found", user.Id)
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
	// Make a get request to the UserAPI to retrieve the user
	userUrl := engine.Configuration.userApiUrl
	r, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/drivers/%s", userUrl, id),
		nil,
	)
	r.Header.Add("Authorization", fmt.Sprintf("Token %s", engine.Configuration.creds.Token))
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
	engine.UserRepo.Add(*user)
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
	// Store the room in the repository.
	engine.RoomRepo.Add(*room)

	go room.Run()
	engine.JoinRoom(id, driver)
	return nil
}

func (engine *ChatEngine) Login() {
	// Login to the UserAPI
	userUrl := engine.Configuration.userApiUrl
	login := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    engine.Configuration.creds.Email,
		Password: engine.Configuration.creds.Password,
	}
	credsBody, _ := json.Marshal(login)
	r, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/login/", userUrl),
		bytes.NewBuffer(credsBody),
	)
	r.Header.Add("Content-Type", "application/json")
	response, err := engine.HttpClient.Do(r)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Fatalf("error: user api login falied: %v", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	token := struct {
		Token string `json:"token"`
	}{}
	if err := json.Unmarshal(body, &token); err != nil {
		log.Fatalf("error: could not initialize: %v", err)
	}
	engine.Configuration.creds.Token = token.Token
}

func (engine *ChatEngine) Initialize() {
	// Clear the database
	engine.UserRepo.Clear()
	engine.RoomRepo.Clear()

	// Request all routes
	routeUrl := engine.Configuration.routeApiUrl
	r, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/routes", routeUrl),
		nil,
	)
	r.Header.Add("Authorization", fmt.Sprintf("Token %s", engine.Configuration.creds.Token))
	response, err := engine.HttpClient.Do(r)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Fatalf("error: could not initialize: %v", ErrUserApiRequestFailed)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	// Parse the routes into objects
	routeList := make([]struct {
		Id         int    `json:"id"`
		Driver     int    `json:"driver"`
		DestAlias  string `json:"destinationAlias"`
		Passengers []struct {
			Id int `json:"id"`
		} `json:"passengers"`
	}, 0)
	if err = json.Unmarshal(body, &routeList); err != nil {
		log.Fatalf("error: could not initialize: %v", ErrRouteUnmarshalFailed)
	}
	// Load the users and open the rooms
	for _, route := range routeList {
		driverId := fmt.Sprintf("%d", route.Driver)
		routeId := fmt.Sprintf("%d", route.Id)

		if err := engine.LoadUser(driverId); err != nil {
			log.Printf("error: could not load user %d: %v", route.Driver, err)
		}
		// Open the room
		if err := engine.OpenRoom(routeId, route.DestAlias, driverId); err != nil {
			log.Printf("error: could not open room %d: %v", route.Id, err)
		}
		// Load the users joined to the route
		for _, user := range route.Passengers {
			userId := fmt.Sprintf("%d", user.Id)
			if err := engine.LoadUser(userId); err != nil {
				log.Printf("error: could not load user %d: %v", user.Id, err)
			}
			// Join the users to the room
			if err := engine.JoinRoom(routeId, userId); err != nil {
				log.Printf("error: could not join user %d to room %d: %v", user.Id, route.Id, err)
			}
		}

	}
}

// NewChatEngine creates a new instance of ChatEngine.
func NewChatEngine(conn *sql.DB, conf *Configuration) *ChatEngine {
	engine := &ChatEngine{
		Configuration: conf,
		HttpClient:    &http.Client{},
		Users:         make(map[string]*User),
		Rooms:         make(map[string]*Room),
		UserRepo:      SqlUserRepository{Db: conn},
		RoomRepo:      SqlRoomRepository{Db: conn},
	}
	engine.Server = NewWsServer(engine)
	return engine
}
