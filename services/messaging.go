package services

import (
	"encoding/json"
	"log"
	"websocket-server/connections"
	"websocket-server/models"

	"github.com/gorilla/websocket"
)

func HandleMessage(senderID string, message []byte) {
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Invalid message from %s: %v\n", senderID, err)
		return
	}

	// Ensure the message has a valid recipient
	if msg.RecipientID == "" {
		log.Printf("Message from %s missing recipient ID\n", senderID)
		return
	}

	conn, ok := connections.GetConnection(msg.RecipientID, "")
	if !ok {
		log.Printf("Recipient %s not connected\n", msg.RecipientID)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("Failed to send message to %s: %v\n", msg.RecipientID, err)
	}
}
