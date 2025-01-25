package routes

import (
	"net/http"
	"websocket-server/handlers"
)

// RegisterUserRoutes sets up routes for the User Service
func RegisterUserRoutes(mux *http.ServeMux) {
	// User-related routes
	mux.HandleFunc("/login", handlers.LoginHandler)       // POST for user login
	mux.HandleFunc("/register", handlers.RegisterHandler) // POST for user registration
	// mux.HandleFunc("/user/logs", handlers.GetUserLogsHandler)       // GET logs for a user
	// mux.HandleFunc("/user/details", handlers.GetUserDetailsHandler) // GET user details
}
