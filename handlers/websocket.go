package handlers

import (
	"log"
	"net/http"
	"websocket-server/authentication"
	"websocket-server/connections"
	"websocket-server/services"

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

	userID, err := authentication.ValidateToken(token)
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
