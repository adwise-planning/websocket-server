package repository

import (
	"database/sql"
	"time"
	"websocket-server/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) LogEvent(event models.Event) error {
	query := `INSERT INTO events (user_id, action, timestamp) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, event.UserID, event.Action, time.Now())
	return err
}

func (r *EventRepository) GetEventsByUserID(userID string) ([]models.Event, error) {
	query := `SELECT user_id, action, timestamp FROM events WHERE user_id = $1 ORDER BY timestamp DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.UserID, &event.Action, &event.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}