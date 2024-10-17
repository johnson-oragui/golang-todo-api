package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/johnson-oragui/golang-todo-api/middleware"
	"github.com/johnson-oragui/golang-todo-api/routes"
)

// myHandler sets the server routes
func myHandler() http.Handler {
	router := mux.NewRouter()

	baseRouter := routes.NewBaseRouter() // Base Handler
	userRouter := routes.NewUserRouter() // Users Handler
	todoRouter := routes.NewTodoRouter() // Todos Handler

	// Define handlers
	router.HandleFunc("/", baseRouter.HomeHandler).Methods("GET")                                                                     // root handler
	router.HandleFunc("/api/v1/about", baseRouter.HandleAboutPage).Methods("GET")                                                     // About page handler
	router.HandleFunc("/api/v1/auth/register", userRouter.HandleRegister)                                                             // POST
	router.HandleFunc("/api/v1/auth/login", userRouter.HandleLogin).Methods("POST")                                                   // POST
	router.Handle("/api/v1/users", middleware.JWTAuthMiddleware(http.HandlerFunc(userRouter.HandleUsers)))                            // GET, PUT, DELETE
	router.Handle("/api/v1/users/todos/{todo_id}", middleware.JWTAuthMiddleware(http.HandlerFunc(todoRouter.HandleTodos)))            // GET, PUT, DELETE
	router.Handle("/api/v1/users/todos", middleware.JWTAuthMiddleware(http.HandlerFunc(todoRouter.HandleGetTodos))).Methods("GET")    // GET
	router.Handle("/api/v1/users/todos", middleware.JWTAuthMiddleware(http.HandlerFunc(todoRouter.HandleCreateTodo))).Methods("POST") // POST
	return middleware.LogginMiddleware(router)
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
