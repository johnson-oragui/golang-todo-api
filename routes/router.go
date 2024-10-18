package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnson-oragui/golang-todo-api/middleware"
)

// myHandler sets the server routes
func MyHandler() http.Handler {
	router := mux.NewRouter()

	baseRouter := NewBaseRouter() // Base Handler
	userRouter := NewUserRouter() // Users Handler
	todoRouter := NewTodoRouter() // Todos Handler

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
