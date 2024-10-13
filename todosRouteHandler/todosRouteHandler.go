package todosRouteHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/johnson-oragui/golang-todo-api/schema"
)

type Router struct{}

func New() *Router {
	return &Router{}
}

func (r *Router) HandleTodos(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.HandleGetTodos(w, req)
	case http.MethodPost:
		r.HandleCreateTodo(w, req)
	case http.MethodPut:
		r.HandleUpdateTodo(w, req)
	case http.MethodDelete:
		r.HandledeleteTodo(w, req)
	default:
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Creates Todo POST /api/v1/users/{username}/todos
func (r *Router) HandleCreateTodo(w http.ResponseWriter, req *http.Request) {
	// check for content-type
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		log.Println("Content-type must be application/json")
		http.Error(w, "Wrong content-type", http.StatusUnsupportedMediaType)
		return
	}

	// get the username from url
	vars := mux.Vars(req)

	username := vars["username"]

	if username == "" {
		log.Println("username is not passed")
		http.Error(w, "username is not passed", http.StatusBadRequest)
		return
	}

	// create a nil TodoSchema struct
	todoInput := schema.TodoSchema{}

	// save the request body to the nill struct
	if err := json.NewDecoder(req.Body).Decode(&todoInput); err != nil {
		log.Println("Error decoding json")
		http.Error(w, "Invalid JSON", http.StatusUnsupportedMediaType)
		return
	}

	// defer closing of request body
	defer req.Body.Close()

	// check if user has a todo entry
	userTodos, exists := schema.TodosDataBase.User[username]
	if !exists {
		// create an empty entry if user has no todo entry
		userTodos = schema.Todos{}
	}

	// add the payload to the list of user todos
	userTodos.AllTodos = append(userTodos.AllTodos, todoInput)

	// save the list to the database
	schema.TodosDataBase.User[username] = userTodos

	// create a response payload
	response := schema.TodoResponse{
		Response: schema.Response{
			Message:    "Todo created successfully",
			StatusCode: 201,
		},
		Data: todoInput,
	}

	// return the payload
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding json")
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
}

func (r *Router) HandleGetTodos(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	username := vars["username"]

	if username == "" {
		log.Println("username is not passed")
		http.Error(w, "username is not passed", http.StatusBadRequest)
		return
	}

	// check if user exists
	_, exists := schema.Database.Users[username]

	if !exists {
		log.Printf("username %v does not exist", username)
		http.Error(w, "username does not exist", http.StatusBadRequest)
		return
	}

	todos := schema.TodosDataBase.User[username].AllTodos
	if !exists {
		log.Printf("username %v does not exist", username)
		http.Error(w, "user does not have a todo entry yet", http.StatusBadRequest)
		return
	}

	response := schema.TodoResponse{
		Response: schema.Response{
			Message:    "Todos retrieved successfully",
			StatusCode: 200,
		},
		Data: todos,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func (r *Router) HandleUpdateTodo(w http.ResponseWriter, req *http.Request) {}

func (r *Router) HandledeleteTodo(w http.ResponseWriter, req *http.Request) {}
