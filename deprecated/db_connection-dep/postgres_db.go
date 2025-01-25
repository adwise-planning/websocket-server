package db_connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

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

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Println("Database connection established.")
}
