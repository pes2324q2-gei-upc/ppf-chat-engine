// main is the entry point of the ppf-chat-engine application.
// It parses command line flags, registers API handlers, and starts the HTTP server.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
)

var addr *string

func main() {
	flag.Parse()
	addr = flag.String("addr", ":8082", "http service address")

	go chat.DefaultChatEngine.Run()
	api.RegisterHandlers()

	log.Fatal(http.ListenAndServe(*addr, nil))
}
