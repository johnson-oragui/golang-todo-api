package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/johnson-oragui/golang-todo-api/baseRouteHandler"
	"github.com/johnson-oragui/golang-todo-api/todosRouteHandler"
	"github.com/johnson-oragui/golang-todo-api/usersRouteHandler"
	"github.com/johnson-oragui/golang-todo-api/utils"
)

// myHandler sets the server routes
func myHandler() http.Handler {
	router := mux.NewRouter()

	baseRouteHandler := baseRouteHandler.New()   // Base Handler
	usersRouteHandler := usersRouteHandler.New() // Users Resource Handler
	todosRouteHandler := todosRouteHandler.New()

	// Define handlers
	router.HandleFunc("/", baseRouteHandler.HomeHandler).Methods("GET")                          // root handler
	router.HandleFunc("/api/v1/about", baseRouteHandler.HandleAboutPage).Methods("GET")          // About page handler
	router.HandleFunc("/api/v1/auth/register", usersRouteHandler.HandleRegister).Methods("POST") // POST
	router.HandleFunc("/api/v1/users/{username}", usersRouteHandler.HandleUsers)                 // GET, PUT, DELETE
	router.HandleFunc("/api/v1/users/{username}/todos", todosRouteHandler.HandleTodos)           // GET, PUT, DELETE
	return utils.LogginMiddleware(router)
}

func main() {
	server := &http.Server{
		Addr:           ":5000",
		Handler:        myHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server running on http://localhost:5000")
	log.Fatal(server.ListenAndServe())
}
