# Simple ToDo API

This is a simple RESTful API for managing users and their to-dos. It uses Go's `net/http` package and stores data in-memory, meaning the data is lost when the server is restarted. The API features user registration, login with JWT-based authentication, and CRUD operations for managing users and their to-do items. 

The API was built using:
- `net/http` for handling HTTP requests
- `gorilla/mux` for routing
- `crypto/bcrypt` for password hashing
- `jwt-go` for JWT token management
- `testing` for unit testing
- `.air.toml` for hot reloading during development

## Endpoints

### Authentication & User Endpoints

1. **POST `/api/v1/auth/register`** - Register a new user
   - **Body**:
     ```json
     {
       "username": "string",
       "password": "string",
       "email": "string",
       "first_name": "string",
       "last_name": "string"
     }
     ```
   - **Response**:
     ```json
     {
       "message": "User Registered successfully",
       "status_code": 201,
       "data": {
         "id": 1234,
         "username": "testuser",
         "first_name": "testuser",
         "last_name": "testuser",
         "email": "testuser@gmail.com"
       }
     }
     ```

2. **POST `/api/v1/auth/login`** - Login and get JWT access token
   - **Body**:
     ```json
     {
       "username": "string",
       "password": "string"
     }
     ```
   - **Response**:
     ```json
     {
       "access_token": "jwt_token"
     }
     ```

3. **GET `/api/v1/users` (Protected)** - Get user details
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Response**:
     ```json
     {
       "id": 1234,
       "username": "testuser",
       "email": "testuser@gmail.com",
       "first_name": "testuser",
       "last_name": "testuser"
     }
     ```

4. **PUT `/api/v1/users` (Protected)** - Update user details
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Body**:
     ```json
     {
       "email": "newemail@gmail.com",
       "first_name": "newFirstName",
       "last_name": "newLastName"
     }
     ```
   - **Response**:
     ```json
     {
       "message": "User updated successfully",
       "status_code": 200,
       "data": {
         "id": 1234,
         "username": "testuser",
         "first_name": "newFirstName",
       "last_name": "newLastName",
         "email": "newemail@gmail.com"
       }
     }
     ```

5. **DELETE `/api/v1/users` (Protected)** - Delete the current user account
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Response**:
     ```json
     {
       "message": "User deleted successfully",
       "status_code": 200
     }
     ```

### To-Do Endpoints (Protected)

1. **GET `/api/v1/users/todos/{todo_id:int}` (Protected)** - Get a specific to-do item
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Response**:
     ```json
     {
       "id": 1,
       "todo": "Buy groceries",
       "completed": false
     }
     ```

2. **GET `/api/v1/users/todos` (Protected)** - Get all to-do items for the user
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Response**:
     ```json
     [
       {
         "id": 1,
         "todo": "Buy groceries",
       "completed": false
       },
       {
         "id": 2,
         "todo": "Walk Dog at 6p.m",
         "completed": false
       }
     ]
     ```

3. **POST `/api/v1/users/todos` (Protected)** - Create a new to-do item
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Body**:
     ```json
     {
       "todo": "string",
       "completed": false | true
     }
     ```
   - **Response**:
     ```json
     {
       "message": "To-do item created successfully",
       "status_code": 201
     }
     ```

4. **PUT `/api/v1/users/todos/{todo_id:int}` (Protected)** - Update a specific to-do item
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Body**:
     ```json
     {
       "todo": "string",
       "completed": false | true
     }
     ```
   - **Response**:
     ```json
     {
       "message": "To-do item updated successfully",
       "status_code": 200
     }
     ```

5. **DELETE `/api/v1/users/todos/{todo_id:int}` (Protected)** - Delete a specific to-do item
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Response**:
     ```json
     {
       "message": "To-do item deleted successfully",
       "status_code": 200
     }
     ```

## How to Run Locally

### Prerequisites
- Go 1.23 or higher
- Optional: `.air.toml` for hot reloading (requires [Air](https://github.com/cosmtrek/air))

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/johnson-oragui/golang-todo-api.git
   cd golang-todo-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. (Optional) Configure Air for hot-reloading. If you have Air installed, you can run:
   ```bash
   air
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

5. The API will be running at `http://localhost:5080`.

### Running Tests

To run unit tests, simply execute the following command:
```bash
go test ./tests -v
```

## Authorization

The API uses JWT (JSON Web Token) for protected routes. After registering, you'll need to login and use the access token provided in the `Authorization` header for subsequent requests.

**Example:**
```bash
Authorization: Bearer <your-access-token>
```

## In-Memory Database

This API uses an in-memory database, which means all data is lost when the server is restarted. The database structure is initialized in the `schema` package.

## Project Structure

```
├── main.go                    # Entry point for the application
├── routes                     # Defines HTTP routes and handlers
├── schema                     # In-memory database and schemas
├── auth                       # JWT and password utilities
├── tests                      # Test cases for API
├── .air.toml                  # Hot reload configuration file
└── go.mod                     # Go module dependencies
```

## Future Improvements
- Implement persistent storage using a database (e.g., PostgreSQL or MongoDB).
- Add token expiration and refresh token mechanism.
- Enhance validation for inputs.
- Add more test coverage for edge cases.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Enjoy building with Go! 🚀