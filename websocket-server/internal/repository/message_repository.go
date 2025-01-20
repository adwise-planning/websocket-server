package repository

import (
	"database/sql"
	"time"
	"websocket-server/internal/models"
)

// MessageRepository provides methods for interacting with message data
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new instance of MessageRepository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// CreateMessage inserts a new message into the database
func (r *MessageRepository) CreateMessage(message *models.Message) error {
	query := `INSERT INTO messages (sender_id, recipient_id, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, message.SenderID, message.RecipientID, message.Content, time.Now())
	return err
}

// GetMessages retrieves messages between two users
func (r *MessageRepository) GetMessages(senderID, recipientID string) ([]models.Message, error) {
	query := `SELECT id, sender_id, recipient_id, content, created_at FROM messages WHERE (sender_id = $1 AND recipient_id = $2) OR (sender_id = $2 AND recipient_id = $1) ORDER BY created_at`
	rows, err := r.db.Query(query, senderID, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.RecipientID, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// DeleteMessage removes a message from the database
func (r *MessageRepository) DeleteMessage(messageID string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, messageID)
	return err
}