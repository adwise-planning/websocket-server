package routes

import (
	"net/http"
	"websocket-server/handlers"
)

// RegisterAuthRoutes sets up routes for the Authentication Service
func RegisterAuthRoutes(mux *http.ServeMux) {
	// Authentication-related routes
	mux.HandleFunc("/login", handlers.LoginHandler) // POST for user login
	// mux.HandleFunc("/validate", handlers.ValidateTokenHandler) // POST for token validation
}
