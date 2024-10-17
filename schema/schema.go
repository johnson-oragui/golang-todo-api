package schema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/johnson-oragui/golang-todo-api/utils"
)

type LoginSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type UserSchemaInput struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserBase struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

type UserSchemaOutput struct {
	Message    string   `json:"message"`
	StatusCode int      `json:"status_code"`
	Data       UserBase `json:"data"`
}

type UsersDataBase struct {
	Users map[string]UserBase
}

type TodoDataBase struct {
	User map[string]Todos
}

type TodoSchema struct {
	ID        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

type TodoResponse struct {
	Response Response
	Data     any `json:"data"`
}
type Todos struct {
	AllTodos []TodoSchema
}

// Simulated global database
var Database UsersDataBase = UsersDataBase{
	Users: map[string]UserBase{},
}

// Simulated global database
var TodosDataBase TodoDataBase = TodoDataBase{
	User: map[string]Todos{},
}

func (u *UserSchemaInput) ValidateUserBase() error {
	// validate username
	if len(strings.TrimSpace(u.Username)) < 3 {
		return fmt.Errorf("username must be atleast 3 characters long, input: '%v'", u.Username)
	}
	// Disallowed characters for username
	notAllowedChars := "!@#$%^&*()_| \\/+?><'\""

	if utils.ContainsAny(u.Username, notAllowedChars) {
		return fmt.Errorf("username cannot contain any of the following characters: %v", notAllowedChars)
	}

	// validate first_name
	if len(strings.TrimSpace(u.FirstName)) < 3 {
		return fmt.Errorf("firstname must be atleast 3 characters long, input: '%v'", u.FirstName)
	}

	// Check for disallowed characters in first_name
	notAllowedChars = "1234567890!@#$%^&*()_| \\/+?><'\""

	if utils.ContainsAny(u.FirstName, notAllowedChars) {
		return fmt.Errorf("firstname cannot contain any of the following characters: %v", notAllowedChars)
	}

	// validate last_name
	if len(strings.TrimSpace(u.LastName)) < 3 {
		return fmt.Errorf("lastname must be atleast 3 characters long, input: '%v'", u.LastName)
	}

	// Check for disallowed characters in last_name
	if utils.ContainsAny(u.LastName, notAllowedChars) {
		return fmt.Errorf("lastname cannot contain any of the following characters: %v", notAllowedChars)
	}

	// validate password
	if len(strings.TrimSpace(u.Password)) < 6 {
		return fmt.Errorf("password must be atleast 6 characters long, input: '%v'", u.Password)
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return err
	}

	// validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,8}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
