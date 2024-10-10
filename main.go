package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data"`
}

// root handler function
func homeHandler(w http.ResponseWriter, req *http.Request) {
	res := Response{
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

// api homehandler
func apiHomeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:    "API Home page!",
		StatusCode: 200,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
	}
}

// myHandler sets the server routes
func myHandler() http.Handler {
	mux := http.NewServeMux()

	// Define handlers
	mux.HandleFunc("/", homeHandler) // root handler
	mux.HandleFunc("/api/v1", apiHomeHandler)
	return logginMiddleware(mux)
}

// loggingMiddleware is a middleware that logs HTTP requests
func logginMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "applicaton/json")
		logger := make(map[string]string)
		header := make(map[string]string)

		logger["Method"] = req.Method
		logger["Path"] = req.URL.Path
		logger["RemoteAddr"] = req.RemoteAddr

		// Convert the header to a string format that can be logged
		for key, values := range req.Header {
			header[key] = fmt.Sprintf("%s", values)
		}

		// Log request details
		log.Printf("{ method: %v, path: %v, User-Agent: %v}", logger["Method"], logger["Path"], header["User-Agent"])

		// Call the next handler in the chain
		next.ServeHTTP(w, req)
	},
	)
}

func main() {
	server := &http.Server{
		Addr:           ":5000",
		Handler:        myHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Server running on http://localhost:5000")
	log.Fatal(server.ListenAndServe())
}
