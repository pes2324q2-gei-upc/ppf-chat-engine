// main is the entry point of the ppf-chat-engine application.
// It parses command line flags, registers API handlers, and starts the HTTP server.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
	_ "github.com/pes2324q2-gei-upc/ppf-chat-engine/docs"
	swag "github.com/swaggo/http-swagger/v2"
)

//	@title		Chat Engine API
//	@BasePath	/
func main() {
	flag.Parse()
	addr := flag.String("addr", ":8082", "http service address")

	ctrl := api.NewChatController(mux.NewRouter(), chat.NewChatEngine())
	http.Handle("/", ctrl.Router)

	ctrl.Router.PathPrefix("/swagger").Handler(swag.Handler(
		swag.URL("http://localhost:8082/swagger/doc.json"),
	)).Methods(http.MethodGet)

	go ctrl.Engine.Server.Run()

	log.Fatal(http.ListenAndServe(*addr, nil))
}
