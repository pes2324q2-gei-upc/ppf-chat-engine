package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

type ChatApiController struct {
	engine *chat.ChatEngine
}

// RegisterHandlers registers the HTTP request handlers.
func (controller *ChatApiController) RegisterHandlers() {
	log.Println("Registering handlers")
	http.HandleFunc("/", controller.RootHandler)
	http.HandleFunc("/connect/<userId>", controller.DefaultConnectHandler)
}

// RootHandler handles the root HTTP request, always sends 200 status code and an empty body.
func (controller *ChatApiController) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// WsHandler handles the WebSocket request.
func (controller *ChatApiController) DefaultConnectHandler(w http.ResponseWriter, r *http.Request) {
	id := chat.UserId(mux.Vars(r)["userId"])
	// If the user does not exist on the engine, load it.
	if !controller.engine.Exists(chat.UserId(id)) {
		if err := controller.engine.LoadUser(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}
	if err := controller.engine.ConnectUser(id, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /user/<userId>/join/<roomId>
// POST /user/<userId>/leave/<roomId>
//
// POST /room
