package routes

import (
	"net/http"
	"websocket-server/handlers"
)

// RegisterMessagingRoutes sets up WebSocket routes
func RegisterMessagingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", loggingMiddleware(handlers.WebSocketHandler)) // WebSocket connection endpoint
	// mux.HandleFunc("/generate-token", loggingMiddleware(handlers.GenerateTokenHandler)) // Token generation endpoint
}
