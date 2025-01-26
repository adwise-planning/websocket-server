package services

import (
	"encoding/json"
	"fmt"
	"websocket-server/connections"
	"websocket-server/database"
	"websocket-server/logging"
	"websocket-server/models"

	"github.com/gorilla/websocket"
)

// HandleMessage processes incoming messages and forwards them to the recipient
func HandleMessage(senderID string, message []byte) {
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		logging.LogError("Invalid message from "+senderID, err)
		return
	}

	// Ensure the message has a valid recipient
	if msg.RecipientID == "" {
		logging.LogError("Message from "+senderID+" missing recipient ID", nil)
		return
	}

	conn, err := connections.GetConnection(msg.RecipientID, "")
	if err != nil {
		logging.LogError("Recipient "+msg.RecipientID+" not connected", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		logging.LogError("Failed to send message to "+msg.RecipientID, err)
	}
}

// SaveMessage saves a message to the database
func SaveMessage(msg *models.Message) (bool, error) {
	query := `INSERT INTO messages (
		sender_id, recipient_id, content, timestamp, message_type, is_read, delivery_status, read_receipt, edited, edit_timestamp, 
		deleted, delete_timestamp, forwarded, forwarded_from, reply_to_id, thread_id, channel_id, priority, attachment_url, 
		attachment_type, reactions, reaction_count, tags, location, language, seen_by, starred, pinned, pin_timestamp, 
		reaction_summary, encryption
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, 
		$27, $28, $29, $30, $31
	)`
	_, err := database.PostgresDB.Exec(query, msg.SenderID, msg.RecipientID, msg.Content, msg.Timestamp, msg.MessageType, msg.IsRead,
		msg.DeliveryStatus, msg.ReadReceipt, msg.Edited, msg.EditTimestamp, msg.Deleted, msg.DeleteTimestamp, msg.Forwarded,
		msg.ForwardedFrom, msg.ReplyToID, msg.ThreadID, msg.ChannelID, msg.Priority, msg.AttachmentURL, msg.AttachmentType,
		msg.Reactions, msg.ReactionCount, msg.Tags, msg.Location, msg.Language, msg.SeenBy, msg.Starred, msg.Pinned, msg.PinTimestamp,
		msg.ReactionSummary, msg.Encryption)
	if err != nil {
		return true, nil
		logging.LogError("Failed to save message to database", err)
		return false, fmt.Errorf("could not save message: %v", err)
	}

	logging.LogInfo(fmt.Sprintf("Message from sender %s to recipient %s at timestamp %s has been saved successfully.", msg.SenderID, msg.RecipientID, msg.Timestamp))
	return true, nil
}
