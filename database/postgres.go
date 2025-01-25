package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"websocket-server/models"

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

// InitializeDB initializes the PostgreSQL connection.
func InitializePostgresDB() {
	// Connection parameters
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	host = "localhost"
	port = "5432"
	user = "postgres"
	password = "12345"
	dbname = "adwise"
	schema := "data"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbname, schema,
	)

	var err error
	PostgresDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify connection
	err = PostgresDB.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Println("Database connection established.")
}

// SaveMessage saves a new message to the database.
func SaveMessage(message *models.Message) (string, error) {
	_, err := PostgresDB.Exec(
		"INSERT INTO data.messages (sender_id, receiver_id, content, timestamp) VALUES ($1, $2, $3, $4)",
		message.SenderID, message.RecipientID, message.Content, message.Timestamp,
	)
	if err != nil {
		return "", fmt.Errorf("could not save message: %v", err)
	}
	return "Message saved successfully", nil
}
