package models

// Message represents a user-to-user text message
type Message struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
	Timestamp   string `json:"timestamp"`
}
