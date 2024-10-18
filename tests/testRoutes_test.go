package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johnson-oragui/golang-todo-api/routes"
	"github.com/johnson-oragui/golang-todo-api/schema"
)

var accessToken string

type TodoOutput struct {
	Message    string            `json:"message"`
	StatusCode int               `json:"status_code"`
	Data       schema.TodoSchema `json:"data"`
}

var registerPayload map[string]string = map[string]string{
	"username":   "testuser",
	"first_name": "testuser",
	"last_name":  "testuser",
	"password":   "Testuser1234#",
	"email":      "testuser@gmail.com",
}

var loginPayload map[string]string = map[string]string{
	"username": "testuser",
	"password": "Testuser1234#",
}

var todoOnePayload map[string]any = map[string]any{
	"todo":      "I am coming home",
	"completed": false,
}

func TestRegister(t *testing.T) {
	router := routes.MyHandler()

	// register user

	payload, err := json.Marshal(registerPayload)
	if err != nil {
		log.Println(err)
		t.Fatal("could not Marshal")
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	responseRecord := httptest.NewRecorder()

	router.ServeHTTP(responseRecord, req)

	if responseRecord.Code != http.StatusCreated {
		t.Fatalf("expected status code 201, but got %v", responseRecord.Code)
	}

	newUser := schema.UserSchemaOutput{}

	err = json.Unmarshal(responseRecord.Body.Bytes(), &newUser)
	if err != nil {
		t.Fatalf("could  not Unmarshal Json, %v", err)
	}

	if newUser.Data.Email != registerPayload["email"] {
		t.Fatalf("expected email to be testuser@gmail.com, but got %v", newUser.Data.Email)
	}
	if newUser.Data.FirstName != registerPayload["first_name"] {
		t.Fatalf("expectec first_name to be testuser, but got %v", newUser.Data.FirstName)
	}
	if newUser.Data.LastName != registerPayload["last_name"] {
		t.Fatalf("expectec last_name to be testuser, but got %v", newUser.Data.LastName)
	}
	if newUser.Data.Username != registerPayload["username"] {
		t.Fatalf("expectec username to be testuser, but got %v", newUser.Data.Username)
	}
}

func TestLogin(t *testing.T) {
	router := routes.MyHandler()

	payload, err := json.Marshal(loginPayload)
	if err != nil {
		log.Println(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		log.Println(responseRecorder.Body.String())
		t.Fatalf("expected code to be 200, but got %v", responseRecorder.Code)
	}

	type LoginResponse struct {
		Message    string            `json:"message"`
		StatusCode int               `json:"status_code"`
		Data       map[string]string `json:"data"`
	}

	loginResponse := LoginResponse{}

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &loginResponse)
	if err != nil {
		log.Println(err)
	}

	accessToken = loginResponse.Data["access_token"]

	if loginResponse.Data["access_token"] == "" {
		t.Fatalf("expected to have access_token, but got %v", loginResponse.Data["access_token"])
	}

}

func TestGetUser(t *testing.T) {
	// setup router
	router := routes.MyHandler()

	bearer := fmt.Sprintf("Bearer %v", accessToken)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected to get 200, but got %v", responseRecorder.Code)
	}

	userResPayload := schema.UserSchemaOutput{}

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &userResPayload)
	if err != nil {
		log.Println(err)
	}

	if userResPayload.Data.Email != "testuser@gmail.com" {
		t.Fatalf("expected to have testuser@gmail.com, but got %v", userResPayload.Data.Email)
	}
	if userResPayload.Data.FirstName != "testuser" {
		t.Fatalf("expectec first_name to be testuser, but got %v", userResPayload.Data.FirstName)
	}
	if userResPayload.Data.LastName != "testuser" {
		t.Fatalf("expectec last_name to be testuser, but got %v", userResPayload.Data.LastName)
	}
	if userResPayload.Data.Username != "testuser" {
		t.Fatalf("expectec username to be testuser, but got %v", userResPayload.Data.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	router := routes.MyHandler()

	userUpdateInput := map[string]string{
		"first_name": "testusergreat",
	}

	payload, err := json.Marshal(userUpdateInput)
	if err != nil {
		log.Println(err)
	}
	req, _ := http.NewRequest("PUT", "/api/v1/users", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != 201 {
		t.Fatalf("expected to get 201, but got %v", responseRecorder.Code)
	}

	userResponse := schema.UserSchemaOutput{}

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &userResponse)
	if err != nil {
		log.Println(err)
	}

	if userResponse.Data.Email != "testuser@gmail.com" {
		t.Fatalf("expected to have testuser@gmail.com, but got %v", userResponse.Data.Email)
	}
	if userResponse.Data.FirstName != "testusergreat" {
		t.Fatalf("expectec first_name to be testuser, but got %v", userResponse.Data.FirstName)
	}
	if userResponse.Data.LastName != "testuser" {
		t.Fatalf("expectec last_name to be testuser, but got %v", userResponse.Data.LastName)
	}
	if userResponse.Data.Username != "testuser" {
		t.Fatalf("expectec username to be testuser, but got %v", userResponse.Data.Username)
	}
}

func TestDeleteUser(t *testing.T) {
	router := routes.MyHandler()

	req, _ := http.NewRequest("DELETE", "/api/v1/users", bytes.NewReader(nil))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	reseponseRecorder := httptest.NewRecorder()

	router.ServeHTTP(reseponseRecorder, req)

	if reseponseRecorder.Code != 202 {
		t.Fatalf("expected to get 202,  but got %v", reseponseRecorder.Code)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// setup router
	router := routes.MyHandler()

	bearer := fmt.Sprintf("Bearer %v", accessToken)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusForbidden {
		t.Fatalf("expected to get 403, but got %v", responseRecorder.Code)
	}

	userResPayload := schema.UserSchemaOutput{}

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &userResPayload)
	if err != nil {
		log.Println(err)
	}

}

func TestGetUserWithoutAuthBearer(t *testing.T) {
	// setup router
	router := routes.MyHandler()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected to get 401, but got %v", responseRecorder.Code)
	}

}

func TestCreateTodo(t *testing.T) {
	router := routes.MyHandler()

	// re-create user
	payload, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// unmarshal register response
	registerResponse := schema.UserSchemaOutput{}

	err := json.Unmarshal(rr.Body.Bytes(), &registerResponse)
	if err != nil {
		log.Println(err)
	}

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected to get 201, but got %v", rr.Code)
	}

	// create todo

	payload, err = json.Marshal(todoOnePayload)
	if err != nil {
		log.Println(err)
	}

	req, _ = http.NewRequest("POST", "/api/v1/users/todos", bytes.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	// confirm status code
	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected to get %v, but got %v", http.StatusCreated, responseRecorder.Code)
	}

	todoResponse := TodoOutput{}

	// unmarshal todo response payload
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &todoResponse)
	if err != nil {
		log.Println(err)
	}

	if todoResponse.Data.Todo != todoOnePayload["todo"] {
		t.Fatalf("expected to have %v, but got %v", todoOnePayload["todo"], todoResponse.Data.Todo)
	}

}

func TestGetTodo(t *testing.T) {
	router := routes.MyHandler()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/todos/1", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected to get %v, but got %v", http.StatusOK, responseRecorder.Code)
	}
	// unmarshal response
	todoResponse := TodoOutput{}

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &todoResponse)
	if err != nil {
		log.Println(err)
	}
	if todoResponse.Data.Todo != todoOnePayload["todo"] {
		t.Fatalf("expected to have %v, but got %v", todoOnePayload["todo"], todoResponse.Data.Todo)
	}
}

func TestGetTodos(t *testing.T) {
	router := routes.MyHandler()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/todos", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected to get %v, but got %v", http.StatusOK, responseRecorder.Code)
	}
	// unmarshal response
	type TodoResponse struct {
		Message    string
		StatusCode int
		Data       []schema.TodoSchema
	}
	todoResponse := TodoResponse{}

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &todoResponse)
	if err != nil {
		log.Println(err)
	}

	todoSlice := []schema.TodoSchema{
		schema.TodoSchema{
			Todo:      todoOnePayload["todo"].(string),
			Completed: todoOnePayload["completed"].(bool),
			ID:        1,
		},
	}
	if len(todoResponse.Data) != len(todoSlice) {
		t.Fatalf("expected to have %v, but got %v", todoOnePayload, todoResponse.Data[0])
	}
}

func TestUpdateTodo(t *testing.T) {
	router := routes.MyHandler()
	newPayload := todoOnePayload
	newPayload["todo"] = "Must be"

	payload, err := json.Marshal(newPayload)
	if err != nil {
		log.Println(err)
	}

	req, _ := http.NewRequest("PUT", "/api/v1/users/todos/1", bytes.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	// confirm status code
	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected to get %v, but got %v", http.StatusCreated, responseRecorder.Code)
	}

	type TodoResponse struct {
		Message    string
		StatusCode int
		Data       []schema.TodoSchema
	}

	todoResponse := TodoResponse{}

	// unmarshal todo response payload
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &todoResponse)
	if err != nil {
		log.Println(err)
	}

	if todoResponse.Data[0].Todo != newPayload["todo"] {
		t.Fatalf("expected to have %v, but got %v", newPayload["todo"], todoResponse.Data[0].Todo)
	}
}

func TestDeleteTodo(t *testing.T) {
	router := routes.MyHandler()

	req, _ := http.NewRequest("DELETE", "/api/v1/users/todos/1", bytes.NewBuffer(nil))
	req.Header.Add("Content-type", "application/json")
	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != 202 {
		t.Fatalf("expected %v but got %v", http.StatusAccepted, responseRecorder.Code)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	router := routes.MyHandler()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/todos/1", bytes.NewBuffer(nil))
	req.Header.Add("Content-Type", "application/json")

	bearer := fmt.Sprintf("Bearer %v", accessToken)
	req.Header.Add("Authorization", bearer)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected to get %v, but got %v", http.StatusNotFound, responseRecorder.Code)
	}
}
