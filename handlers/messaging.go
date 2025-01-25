package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"websocket-server/connections"
	"websocket-server/models"
	"websocket-server/services"
	"websocket-server/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := utils.ValidateToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer conn.Close()

	connections.AddConnection(userID, token, conn)
	defer connections.RemoveConnection(userID, token)

	log.Printf("User %s connected\n", userID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for user %s: %v\n", userID, err)
			break
		}

		services.HandleMessage(userID, message)
	}
}

// GenerateTokenHandler handles token generation for a given user
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query params
	var requestData = models.GenerateTokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := requestData.Username
	device_id := requestData.DeviceID
	email := requestData.Email

	// Generate a token for the user
	token, err := utils.GenerateToken(username, device_id, email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate access token: %v", err), http.StatusInternalServerError)
		return
	}
	// refreshToken
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate refresh token: %v", err), http.StatusInternalServerError)
		return
	}
	// Send the token as a JSON response
	response := models.TokenResponse{
		UserName:     username,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
