package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Required for PostgreSQL driver
)

// DB holds the global database connection instance.
// This connection is used across all database operations in the application.
// It is initialized during application startup via InitDB().
var DB *sql.DB

// InitDB initializes and verifies the database connection.
// It performs the following operations:
//  1. Constructs connection string from environment variables
//  2. Opens database connection
//  3. Verifies connection with ping
//  4. Stores connection in global DB variable
//
// Required environment variables:
//   - DB_HOST: Database server hostname
//   - DB_PORT: Database server port
//   - DB_USER: Database username
//   - DB_PASSWORD: Database password
//   - DB_NAME: Database name
//
// Returns:
//   - error: Any error encountered during connection initialization
func InitDB() error {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Verify connection is active and accessible
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Successfully connected to database")
	return nil
}
