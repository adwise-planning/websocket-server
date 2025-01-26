package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"websocket-server/logging"

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

// InitializePostgresDB initializes the PostgreSQL connection.
func InitializePostgresDB() error {
	// Retrieve connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	schema := os.Getenv("DB_SCHEMA")

	host = "dpg-cu796edds78s73aq6iu0-a.oregon-postgres.render.com"
	port = "5432"
	user = "admin"
	password = "rsg80kYOY6qbbbBfSwWwdcjxF6gYevFP"
	dbname = "adwise"
	schema = "data"

	// Validate environment variables
	if host == "" || port == "" || user == "" || password == "" || dbname == "" || schema == "" {
		logging.LogError("Database connection parameters are not set in environment variables", nil)
		return fmt.Errorf("database connection parameters are not set in environment variables")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s",
		host, port, user, password, dbname, schema,
	)

	var err error
	PostgresDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		logging.LogError("Failed to connect to the database", err)
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Configure connection pooling
	PostgresDB.SetMaxOpenConns(25)
	PostgresDB.SetMaxIdleConns(25)
	PostgresDB.SetConnMaxLifetime(15 * time.Minute)

	// Verify connection
	err = PostgresDB.Ping()
	if err != nil {
		logging.LogError("Could not ping database", err)
		return fmt.Errorf("could not ping database: %v", err)
	}

	logging.LogInfo("Database connection established.")
	return nil
}
