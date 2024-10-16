package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

// loggingMiddleware is a middleware that logs HTTP requests
func LogginMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := make(map[string]string)
		header := make(map[string]any)

		logger["Method"] = req.Method
		logger["Path"] = req.URL.Path
		logger["RemoteAddr"] = req.RemoteAddr

		// Convert the header to a string format that can be logged
		for key, values := range req.Header {
			header[key] = fmt.Sprintf("%s", values)
		}

		// Read the body into memory for logging, then restore it for further processing
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error in logger: %v", err)
		}

		// Log request details
		log.Printf("{ method: %v, path: %v, Body: %v}", logger["Method"], logger["Path"], string(bodyBytes))

		// Restore the body so it can be read again by the actual handler
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Call the next handler in the chain
		next.ServeHTTP(w, req)
	},
	)
}
