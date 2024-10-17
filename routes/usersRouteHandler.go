package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/johnson-oragui/golang-todo-api/auth"
	"github.com/johnson-oragui/golang-todo-api/schema"
	"github.com/johnson-oragui/golang-todo-api/utils"
)

type UserRouter struct{}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

// Users Rosource Route Handler
func (b *UserRouter) HandleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:
		b.HandleGetUser(w, req)
	case http.MethodPut:
		b.HandleUpdateuser(w, req)
	case http.MethodDelete:
		b.HandleDeleteUser(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

// create user handler POST /users
func (s *UserRouter) HandleRegister(w http.ResponseWriter, req *http.Request) {
	var newUser schema.UserSchemaInput
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

	if err := newUser.ValidateUserBase(); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	// check if user already exists
	userExists, exists := schema.Database.Users[newUser.Username]
	if exists {
		log.Println("User already exists, user:", userExists)
		http.Error(w, "User already exists", http.StatusForbidden)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		log.Println("error hashing password")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Construct response data
	data := schema.UserBase{
		Username:  newUser.Username,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		ID:        rand.Intn(10000),
		Password:  hashedPassword,
	}

	// save user to database
	schema.Database.Users[newUser.Username] = data

	res := schema.UserSchemaOutput{
		Message:    "User Registered successfully",
		StatusCode: 200,
		Data:       data,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}
}

func (s *UserRouter) HandleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		log.Println("content-type must be application/json")
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	loginSchema := schema.LoginSchema{}

	if err := json.NewDecoder(req.Body).Decode(&loginSchema); err != nil {
		log.Println("Error Decoding JSON")
		http.Error(w, "Error Decoding JSON", http.StatusInternalServerError)
		return
	}

	notAllowedChars := "!@#$%^&*()_| \\/+?><'\""
	if err := utils.ContainsAny(loginSchema.Username, notAllowedChars); err {
		log.Printf("username must not contain %v", notAllowedChars)
		message := fmt.Sprintf("username must not contain %v", notAllowedChars)
		http.Error(w, message, http.StatusNotFound)
		return
	}

	if err := utils.ValidatePassword(loginSchema.Password); err != nil {
		log.Printf("error: %v", err)
		http.Error(w, fmt.Sprint(err), http.StatusNotFound)
		return
	}

	// check if user exists
	user, exists := schema.Database.Users[loginSchema.Username]
	if !exists {
		log.Printf("user does not exist")
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	err := auth.ComparePasswords(loginSchema.Password, user.Password)

	if err != nil {
		log.Printf("invalid username or password")
		http.Error(w, "invalid username or password", http.StatusForbidden)
		return
	}

	accessToken, err := auth.GenerateJWT(loginSchema.Username)
	if err != nil {
		log.Println(fmt.Sprintln(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := schema.TodoResponse{
		Response: schema.Response{
			StatusCode: 200,
			Message:    "Login Success",
		},
		Data: map[string]string{
			"access_token": accessToken,
		},
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

// fetch user handler GET /users
func (s *UserRouter) HandleGetUser(w http.ResponseWriter, req *http.Request) {
	username, ok := req.Context().Value("username").(string)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusBadRequest)
		return
	}

	user, exists := schema.Database.Users[username]

	if !exists {
		log.Printf("username %s does not exists in the database", username)
		http.Error(w, "User does not exists", http.StatusForbidden)
		return
	}

	res := schema.UserSchemaOutput{
		Message:    "Retrieved successfully",
		StatusCode: 200,
		Data:       user,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error ENcoding JSON", http.StatusInternalServerError)
	}
}

// update user handler PUT /users
func (r *UserRouter) HandleUpdateuser(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	username, ok := req.Context().Value("username").(string)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusBadRequest)
		return
	}

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
	user, exists := schema.Database.Users[username]

	if !exists {
		http.Error(w, "User not Found", http.StatusNotFound)
		return
	}
	notAllowedChars := "1234567890!@#$%^&*()_| \\/+?><'\""

	// update the user
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	if updateUser.FirstName != "" {
		if err := utils.ContainsAny(updateUser.FirstName, notAllowedChars); err {
			log.Printf("firstname must not contain %v", notAllowedChars)
			message := fmt.Sprintf("firstname must not contain %v", notAllowedChars)
			http.Error(w, message, http.StatusNotFound)
			return
		}
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		if err := utils.ContainsAny(updateUser.Username, notAllowedChars); err {
			log.Printf("lastname must not contain %v", notAllowedChars)
			message := fmt.Sprintf("lastname must not contain %v", notAllowedChars)
			http.Error(w, message, http.StatusNotFound)
			return
		}
		user.LastName = updateUser.LastName
	}
	if updateUser.Password != "" {
		if err := utils.ValidatePassword(updateUser.Password); err != nil {
			log.Printf("error: %v", err)
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}
		user.Password = updateUser.Password
	}

	delete(schema.Database.Users, username)

	schema.Database.Users[user.Username] = user

	response := schema.UserSchemaOutput{
		Message:    "Updated successfully",
		StatusCode: 200,
		Data:       schema.Database.Users[user.Username],
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

}

// delete user handler DELETE /users
func (r *UserRouter) HandleDeleteUser(w http.ResponseWriter, req *http.Request) {

	username, ok := req.Context().Value("username").(string)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusBadRequest)
		return
	}

	if _, exists := schema.Database.Users[username]; !exists {
		log.Printf("username %v does not exist", username)
		message := fmt.Sprintf("User %v does not exist", username)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	delete(schema.Database.Users, username)

	response := schema.Response{
		Message:    "User deleted successfully",
		StatusCode: 200,
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("An error occured: %v", err)
		http.Error(w, "Could not ENcode Json response", http.StatusInternalServerError)
	}
}
