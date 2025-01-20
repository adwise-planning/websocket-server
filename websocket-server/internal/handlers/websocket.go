package handlers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"websocket-server/internal/logging"
	"websocket-server/internal/models"
	"websocket-server/internal/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	MessageService *services.MessageService
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.LogError("Error while upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			logging.LogError("Error while reading message: %v", err)
			break
		}

		err = h.MessageService.SaveMessage(msg)
		if err != nil {
			logging.LogError("Error while saving message: %v", err)
			break
		}

		err = conn.WriteJSON(msg)
		if err != nil {
			logging.LogError("Error while writing message: %v", err)
			break
		}
	}
}

func (h *WebSocketHandler) BroadcastMessage(msg models.Message) {
	// Logic to broadcast message to all connected clients
	fmt.Println("Broadcasting message:", msg)
}