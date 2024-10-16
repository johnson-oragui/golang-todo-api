package todosRouteHandler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

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
		r.HandleGetTodo(w, req)
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

	// check if user exists in the users database
	_, exists := schema.Database.Users[username]
	if !exists {
		log.Println("user does not exists in the database", username)
		http.Error(w, "user does not exists in the database", http.StatusForbidden)
		return
	}

	// create a nil TodoSchema struct
	todoInput := schema.TodoSchema{
		ID: rand.Intn(10000),
	}

	// save the request body to the nill struct
	if err := json.NewDecoder(req.Body).Decode(&todoInput); err != nil {
		log.Println("Error decoding json", todoInput)
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

	w.Header().Add("Content-Type", "applicaton/json")
	// return the payload
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding json")
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
}

// Fetch all Todos GET /api/v1/users/{username}/todos
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

	todos, exists := schema.TodosDataBase.User[username]
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
		Data: todos.AllTodos,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// fetch a single Todo GET /api/v1/users/{username}/todos/{todo_id}
func (r *Router) HandleGetTodo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	username := vars["username"]
	todoIdStr := vars["todo_id"]

	if username == "" || todoIdStr == "" {
		log.Println("username or todo_id is not passed")
		http.Error(w, "username or todo_id is not passed", http.StatusBadRequest)
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

	todoId, err := strconv.Atoi(todoIdStr)
	if err != nil {
		log.Println("Invalid todo_id, must be an integer")
		http.Error(w, "Invalid todo_id", http.StatusBadRequest)
		return
	}

	var thatTodo schema.TodoSchema

	for idx, td := range todos {
		if td.ID == todoId {
			thatTodo.ID = todoId
			thatTodo.Completed = todos[idx].Completed
			thatTodo.Todo = todos[idx].Todo
			break
		}
	}

	if thatTodo.ID == 0 && thatTodo.Todo == "" && !thatTodo.Completed {
		log.Println("todo not found")
		http.Error(w, "todo not found", http.StatusInternalServerError)
		return
	}

	response := schema.TodoResponse{
		Response: schema.Response{
			Message:    "Todos retrieved successfully",
			StatusCode: 200,
		},
		Data: thatTodo,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// update a Todo PUT /api/v1/users/{username}/todos/{todo_id}
func (r *Router) HandleUpdateTodo(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		log.Println("Content-type must be application/json")
		http.Error(w, "Wrong content-type", http.StatusUnsupportedMediaType)
		return
	}

	// get username fromparams
	vars := mux.Vars(req)

	username := vars["username"]

	if username == "" {
		log.Println("username not provided")
		http.Error(w, "username not provided", http.StatusBadRequest)
		return
	}

	// decode payload
	todoInput := schema.TodoSchema{}

	if err := json.NewDecoder(req.Body).Decode(&todoInput); err != nil {
		log.Println("invalid JSON")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userTodos, exists := schema.TodosDataBase.User[username]

	if !exists {
		log.Println("invalid JSON")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updated := false
	for i, ts := range userTodos.AllTodos {
		if ts.ID == todoInput.ID {
			if todoInput.Completed {
				userTodos.AllTodos[i].Completed = todoInput.Completed
			}

			if todoInput.Completed {
				userTodos.AllTodos[i].Todo = todoInput.Todo
			}
			updated = true
			break
		}
	}

	if !updated {
		log.Println("Todo not found")
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	schema.TodosDataBase.User[username] = userTodos

	response := schema.TodoResponse{
		Response: schema.Response{
			Message:    "Todos retrieved successfully",
			StatusCode: 200,
		},
		Data: userTodos.AllTodos,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// delete a single Todo DELETE /api/v1/users/{username}/todos/{todo_id}
func (r *Router) HandledeleteTodo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	username := vars["username"]
	todoIdStr := vars["todo_id"]

	if username == "" || todoIdStr == "" {
		log.Println("username or todo_id not provided")
		http.Error(w, "username or todo_id not provided", http.StatusBadRequest)
		return
	}

	todoId, err := strconv.Atoi(todoIdStr)
	if err != nil {
		log.Println("Invalid todo_id, must be an integer")
		http.Error(w, "Invalid todo_id", http.StatusBadRequest)
		return
	}

	userTodos, exists := schema.TodosDataBase.User[username]

	if !exists {
		log.Println("username does not exist")
		http.Error(w, "username does not exist", http.StatusBadRequest)
		return
	}

	for idx, td := range userTodos.AllTodos {
		if td.ID == todoId {
			userTodos.AllTodos = append(userTodos.AllTodos[:idx], userTodos.AllTodos[idx+1:]...)
			break
		}

	}

	schema.TodosDataBase.User[username] = userTodos

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusAccepted)

}
