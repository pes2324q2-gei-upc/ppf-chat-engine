// main is the entry point of the ppf-chat-engine application.
// It parses command line flags, registers API handlers, and starts the HTTP server.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
	"github.com/pes2324q2-gei-upc/ppf-chat-engine/chat"
	_ "github.com/pes2324q2-gei-upc/ppf-chat-engine/docs"
	swag "github.com/swaggo/http-swagger/v2"
)

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
}

// @title		Chat Engine API
// @BasePath	/
func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")
	dbPath := flag.String("db", "chat.db", "database path")
	flag.Parse()

	debug := getEnv("DEBUG", "false") == "true"
	routeApiUrl := getEnv("ROUTE_API_URL", "http://localhost:8080")
	userApiUrl := getEnv("USER_API_URL", "http://localhost:8081")
	mail := getEnv("PPF_MAIL", "admin@ppf.com")
	pass := getEnv("PPF_PASS", "chatengine")

	conf := chat.NewConfiguration(debug, userApiUrl, routeApiUrl, mail, pass)

	db := chat.InitDB("sqlite3", *dbPath)
	router := mux.NewRouter()
	engine := chat.NewChatEngine(db, conf)
	ctrl := api.NewChatController(router, engine)

	engine.Login()
	engine.Initialize()

	// Swagger documentation route
	ctrl.Router.PathPrefix("/swagger").Handler(swag.Handler(
		swag.URL("http://localhost:8082/swagger/doc.json"),
	)).Methods(http.MethodGet)

	go ctrl.Engine.Server.Run()

	http.Handle("/", ctrl.Router)
	log.Printf("info: starting server on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
