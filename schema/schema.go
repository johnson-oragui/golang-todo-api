package schema

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type UserSchemaInput struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Password  string
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
	Response Response
	Data     UserBase `json:"data"`
}

type TodoSchema struct {
	Todo   string `json:"todo"`
	Status string `json:"status"`
}

type TodoResponse struct {
	Response Response
	Data     TodoSchema `json:"data"`
}

type UsersDataBase struct {
	Users map[string]UserBase
}

type TodoDataBase struct {
	Username string
	Todos    map[string]TodoSchema
}
