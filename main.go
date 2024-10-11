package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/johnson-oragui/golang-todo-api/routeHandler"
	"github.com/johnson-oragui/golang-todo-api/utils"
)

// myHandler sets the server routes
func myHandler() http.Handler {
	router := http.NewServeMux()
	
	routeHandler := routeHandler.New()

	// Define handlers
	router.HandleFunc("/", routeHandler.HomeHandler) // root handler
	router.HandleFunc("/api/v1/auth/register", routeHandler.HandleUsers)
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
