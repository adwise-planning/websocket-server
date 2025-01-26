package models

import (
	"time"
)

// Message represents a user-to-user text message
type Message struct {
	SenderID        string     `json:"sender_id"`
	RecipientID     string     `json:"recipient_id"`
	Content         string     `json:"content"`
	Timestamp       time.Time  `json:"timestamp"`
	MessageType     string     `json:"message_type"`     // Type of message (e.g., text, image, video)
	IsRead          bool       `json:"is_read"`          // Indicates if the message has been read
	DeliveryStatus  string     `json:"delivery_status"`  // Status of message delivery (e.g., sent, delivered, failed)
	ReadReceipt     bool       `json:"read_receipt"`     // Indicates if read receipt is enabled
	Edited          bool       `json:"edited"`           // Indicates if the message has been edited
	EditTimestamp   *time.Time `json:"edit_timestamp"`   // Timestamp of the last edit
	Deleted         bool       `json:"deleted"`          // Indicates if the message has been deleted
	DeleteTimestamp *time.Time `json:"delete_timestamp"` // Timestamp of the deletion
	Forwarded       bool       `json:"forwarded"`        // Indicates if the message has been forwarded
	ForwardedFrom   string     `json:"forwarded_from"`   // ID of the original sender if forwarded
	ReplyToID       string     `json:"reply_to_id"`      // ID of the message being replied to
	ThreadID        string     `json:"thread_id"`        // ID of the message thread
	ChannelID       string     `json:"channel_id"`       // ID of the channel where the message was sent
	Priority        int        `json:"priority"`         // Priority level of the message
	AttachmentURL   string     `json:"attachment_url"`   // URL to any attachment if present
	AttachmentType  string     `json:"attachment_type"`  // Type of attachment (e.g., image, video, file)
	Reactions       string     `json:"reactions"`        // Reactions to the message (e.g., like, love, etc.)
	ReactionCount   int        `json:"reaction_count"`   // Total number of reactions
	Tags            string     `json:"tags"`             // Tags associated with the message
	Location        string     `json:"location"`         // Location information if shared
	Language        string     `json:"language"`         // Language of the message content
	SeenBy          string     `json:"seen_by"`          // List of user IDs who have seen the message
	Starred         bool       `json:"starred"`          // Indicates if the message is starred
	Pinned          bool       `json:"pinned"`           // Indicates if the message is pinned
	PinTimestamp    *time.Time `json:"pin_timestamp"`    // Timestamp of when the message was pinned
	ReactionSummary string     `json:"reaction_summary"` // Summary of reactions (e.g., {"like": "5", "love": "3"})
	Encryption      bool       `json:"encryption"`       // Indicates if the message is encrypted
}
