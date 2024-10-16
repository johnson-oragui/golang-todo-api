package baseRouteHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/johnson-oragui/golang-todo-api/schema"
)

type BaseRouter struct{}

func New() *BaseRouter {
	return &BaseRouter{}
}

func (b *BaseRouter) HandleRoute(w http.ResponseWriter, req *http.Request) {
	switch {
	case (req.Method == http.MethodGet && req.URL.Path == "/api/v1/about"):
		b.HandleAboutPage(w, req)
	default:
		b.HomeHandler(w, req)
	}
}

// root handler function  GET
func (s *BaseRouter) HomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	res := schema.Response{
		Message:    "Welcome to the golang todo page!",
		StatusCode: 200,
	}
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
	}
}

// About page handler GET
func (b *BaseRouter) HandleAboutPage(w http.ResponseWriter, req *http.Request) {
	response := schema.Response{
		StatusCode: 200,
		Message:    "This is the about page for the golang todo API",
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding data for about page")
		http.Error(w, "Error encoding data for about page", http.StatusInternalServerError)
	}
}
