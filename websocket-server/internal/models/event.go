package models

import "time"

// Event represents a user event within the application
type Event struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}