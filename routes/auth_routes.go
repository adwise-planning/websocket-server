package routes

import (
	"net/http"
	"websocket-server/handlers"
	"websocket-server/logging"
)

// RegisterAuthRoutes sets up routes for the Authentication Service
func RegisterAuthRoutes(mux *http.ServeMux) {
	// Authentication-related routes
	mux.HandleFunc("/login", loggingMiddleware(handlers.LoginHandler)) // POST for user login
	// mux.HandleFunc("/validate", loggingMiddleware(handlers.ValidateTokenHandler)) // POST for token validation
}

// loggingMiddleware is a middleware that logs the details of each request
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.LogInfo("Received request: " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
		logging.LogInfo("Completed request: " + r.Method + " " + r.URL.Path)
	}
}
