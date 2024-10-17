package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/johnson-oragui/golang-todo-api/auth"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// get token from authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			log.Println("Authorization token not provided")
			http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
			return
		}

		// extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		username, err := auth.DecodeJWT(token)
		if err != nil {
			log.Println("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// add the username to request context and call next handler
		ctx := req.Context()
		ctx = context.WithValue(ctx, "username", username)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
