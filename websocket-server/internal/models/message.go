package models

import "time"

// Message represents a chat message in the system
type Message struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}