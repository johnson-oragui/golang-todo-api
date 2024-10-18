package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johnson-oragui/golang-todo-api/routes"
)

func main() {
	server := &http.Server{
		Addr:           ":5000",
		Handler:        routes.MyHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server running on http://localhost:5000")
	log.Fatal(server.ListenAndServe())
}
