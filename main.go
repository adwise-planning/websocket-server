package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"websocket-server/authentication"
	"websocket-server/db_connection"
	"websocket-server/handlers"
)

// Response structure for token generation
type TokenResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// GenerateTokenHandler handles token generation for a given user
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query params
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// Generate a token for the user
	token, err := authentication.GenerateToken(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the token as a JSON response
	response := TokenResponse{
		UserID: userID,
		Token:  token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {

	db_connection.InitializePostgresDB()

	// WebSocket endpoint
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	// Token generation endpoint
	http.HandleFunc("/generate-token", GenerateTokenHandler)

	port := "10000"
	fmt.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
