package services

import (
	"database/sql"
	"time"
	"websocket-server/internal/models"
)

type EventService struct {
	db *sql.DB
}

func NewEventService(db *sql.DB) *EventService {
	return &EventService{db: db}
}

func (es *EventService) LogEvent(userID string, action string) error {
	event := models.Event{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	query := `INSERT INTO events (user_id, action, timestamp) VALUES ($1, $2, $3)`
	_, err := es.db.Exec(query, event.UserID, event.Action, event.Timestamp)
	return err
}

func (es *EventService) GetUserEvents(userID string) ([]models.Event, error) {
	var events []models.Event

	query := `SELECT user_id, action, timestamp FROM events WHERE user_id = $1 ORDER BY timestamp DESC`
	rows, err := es.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.UserID, &event.Action, &event.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}