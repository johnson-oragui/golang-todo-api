package usersRouteHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnson-oragui/golang-todo-api/schema"
)

// Simulated global database
var Database schema.UsersDataBase = schema.UsersDataBase{
	Users: map[string]schema.UserBase{},
}

// Structure for route handlers
type RouteHandler struct{}

func New() *RouteHandler {
	return &RouteHandler{}
}

// Users Rosource Route Handler
func (r *RouteHandler) HandleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:
		r.HandleGetUser(w, req)
	case http.MethodPut:
		r.HandleUpdateuser(w, req)
	case http.MethodDelete:
		r.HandleDeleteUser(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// create user handler POST /users
func (s *RouteHandler) HandleRegister(w http.ResponseWriter, req *http.Request) {
	var newUser schema.UserBase
	if req.Method != http.MethodPost {
		log.Println("Method Not allowed in register route")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	// Limit the size of the request body
	req.Body = http.MaxBytesReader(w, req.Body, 1048576) // 1MB limit

	// Decode the JSON request body directly into the struct
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	// save user to database
	newUser.ID = 1
	Database.Users[newUser.Username] = newUser

	// Construct response data
	data := schema.UserBase{
		Username:  newUser.Username,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		ID:        1,
	}

	res := schema.UserSchemaOutput{
		Response: schema.Response{
			Message:    "Retrieved successfully",
			StatusCode: 200,
		},
		Data: data,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}
}

// fetch user handler GET /users
func (s *RouteHandler) HandleGetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	if username == "" {
		log.Println("username is missing")
		http.Error(w, "username is missing in the query params", http.StatusBadRequest)
	}

	user, exists := Database.Users[username]

	if !exists {
		log.Printf("username %s does not exists in the database", username)
		http.Error(w, "User does not exists", http.StatusForbidden)
		return
	}

	res := schema.UserSchemaOutput{
		Response: schema.Response{
			Message:    "User retrived successfully",
			StatusCode: 200,
		},
		Data: user,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error ENcoding JSON", http.StatusInternalServerError)
	}
}

// update user handler PUT /users
func (r *RouteHandler) HandleUpdateuser(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	vars := mux.Vars(req)
	// get the user from query params
	username := vars["username"]

	// Decode the request body into the UserSchemaInput struct
	updateUser := schema.UserSchemaInput{}

	err := json.NewDecoder(req.Body).Decode(&updateUser)

	if err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// close
	defer req.Body.Close()

	// retrieve the user from database using the username
	user, exists := Database.Users[username]

	if !exists {
		http.Error(w, "User not Found", http.StatusNotFound)
		return
	}

	// update the user
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Password != "" {
		user.Password = updateUser.Password
	}

	// save the updated user to database
	Database.Users[user.Username] = user

	// Construct the response data
	data := schema.UserBase{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ID:        user.ID,
	}
	response := schema.UserSchemaOutput{
		Response: schema.Response{
			StatusCode: 200,
			Message:    "Updated successfully",
		},
		Data: data,
	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

}

// delete user handler DELETE /users
func (r *RouteHandler) HandleDeleteUser(w http.ResponseWriter, req *http.Request) {

}
