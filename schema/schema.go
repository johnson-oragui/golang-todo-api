package schema

type Response struct {
	Message string `json:"message"`
	StatusCode int `json:"status_code"`
}

type UserSchemaInput struct {
	Username string
	Email string
	Password string
}

type UserBase struct {
	ID int			`json:"id"`
	Username string	`json:"username"`
	FirstName string	`json:"first_name"`
	LastName string	`json:"last_name"`
	Email string	`json:"email"`
	Password string	`json:"password"`
}

type UserSchemaOutput struct {
	Response Response
	Data UserBase `json:"data"`
}

type DataBase struct {
	Users map[string]UserBase
}
