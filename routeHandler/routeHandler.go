package routeHandler

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/johnson-oragui/golang-todo-api/schema"

)

// Structure for route handlers
type RouteHandler struct {}

func New() *RouteHandler {
	return &RouteHandler{}
}



// root handler function
func (s *RouteHandler) HomeHandler(w http.ResponseWriter, req *http.Request) {
	res := schema.Response{
		Message:    "Welcome to the golang todo page!",
		StatusCode: 200,
	}
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
	}
}

// Users Rosource Route Handler
func (r *RouteHandler) HandleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	
	case http.MethodPost:
		r.HandleCreateUsers(w, req)
	case http.MethodGet:
		r.HandleGetUser(w, req)
	}
}

// create user handler
func (s *RouteHandler) HandleCreateUsers(w http.ResponseWriter, req * http.Request) {
	var newUser schema.UserBase

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
	}
	// Limit the size of the request body
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)  // 1MB limit
	
	// Decode the JSON request body directly into the struct
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// save user to database
	Database := schema.DataBase{
		Users: map[string]schema.UserBase{
			newUser.Username: newUser,
		},
	}
	log.Printf("%v", Database)

	// Construct response data
	data := schema.UserBase{
		Username: newUser.Username,
		Email: newUser.Email,
	}
	
	res := schema.UserSchemaOutput{
		Response: schema.Response{
			Message: "Retrieved successfully",
			StatusCode: 200,
		},
		Data: data,
	}
	
	
	
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
}

func (s *RouteHandler) HandleGetUser(w http.ResponseWriter, req *http.Request) {
	Database := &schema.DataBase{}
	users := Database.Users["jay"]
	
	res := schema.UserSchemaOutput{
		Response: schema.Response{
			Message: "User retrived successfully",
			StatusCode: 200,
		},
		Data: users,
	}
	
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error ENcoding JSON", http.StatusInternalServerError)
	}
}