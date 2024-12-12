package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// checkAndCreateDatabase verifies if the required database exists and creates it if not.
// It establishes a connection to PostgreSQL and performs the necessary checks.
//
// Returns:
//   - error: Any error encountered during database verification or creation
func checkAndCreateDatabase() (err error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Printf("DB_HOST: %s", dbHost)
	log.Printf("DB_PORT: %s", dbPort)
	log.Printf("DB_USER: %s", dbUser)
	log.Printf("DB_PASSWORD: %s", dbPassword)
	log.Printf("DB_NAME: %s", dbName)

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %v", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing database connection: %v", closeErr)
		}
	}()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %v", err)
	}

	if !exists {
		log.Printf("Creating database %s...", dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("error creating database: %v", err)
		}
		log.Printf("Database %s created successfully", dbName)
	} else {
		log.Printf("Database %s already exists", dbName)
	}

	return nil
}

// checkAndCreateTables verifies if required tables exist and creates them if not.
// It creates the room_bookings table with necessary indexes for efficient querying.
//
// Parameters:
//   - db *sql.DB: Active database connection
//
// Returns:
//   - error: Any error encountered during table verification or creation
func checkAndCreateTables(db *sql.DB) error {
	var exists bool
	query := `
       SELECT EXISTS (
           SELECT 1 
           FROM information_schema.tables 
           WHERE table_name = 'room_bookings'
       )`
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking table existence: %v", err)
	}

	if !exists {
		log.Println("Creating room_bookings table...")
		query := `
       CREATE TABLE room_bookings (
           id SERIAL PRIMARY KEY,
           room_id VARCHAR(50) NOT NULL,
           date DATE NOT NULL,
           is_booked BOOLEAN NOT NULL,
           rate DECIMAL(10,2) NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           UNIQUE(room_id, date)
       );
       CREATE INDEX idx_room_bookings_room_id ON room_bookings(room_id);
       CREATE INDEX idx_room_bookings_date ON room_bookings(date);
       `

		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error creating tables: %v", err)
		}
		log.Println("Table created successfully")
	} else {
		log.Println("Table room_bookings already exists")
	}

	return nil
}

// generateRoomID creates a random room identifier.
// The ID format is a single uppercase letter followed by three digits (e.g., "A123").
//
// Returns:
//   - string: Generated room identifier
func generateRoomID() string {
	letter := string(rune('A' + rand.Intn(26)))
	number := rand.Intn(900) + 100
	return fmt.Sprintf("%s%d", letter, number)
}

// generateDates creates a sequence of dates starting from today.
// The sequence spans the next 7 months, with one entry per day.
//
// Returns:
//   - []time.Time: Slice of dates from today to 7 months ahead
func generateDates() []time.Time {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 7, 0)

	var dates []time.Time
	currentDate := startDate

	for currentDate.Before(endDate) {
		dates = append(dates, currentDate)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dates
}

// generateMockData creates and inserts mock booking data for rooms.
// It generates random room IDs, sets varying rates and booking status,
// and inserts this data into the database for a 7-month period.
//
// Parameters:
//   - db *sql.DB: Active database connection
//
// Returns:
//   - []string: Slice of generated room identifiers
func generateMockData(db *sql.DB) []string {
	log.Println("Starting mock data generation...")

	var roomIDs []string
	for i := 0; i < 10; i++ {
		roomIDs = append(roomIDs, generateRoomID())
	}

	dates := generateDates()

	for _, roomID := range roomIDs {
		baseRate := 80.0 + rand.Float64()*120.0

		for _, date := range dates {
			rateVariation := 0.8 + rand.Float64()*0.4
			dailyRate := round(baseRate * rateVariation)
			isBooked := rand.Float64() < 0.6

			query := `
           INSERT INTO room_bookings (room_id, date, is_booked, rate)
           VALUES ($1, $2, $3, $4)
           ON CONFLICT (room_id, date) DO NOTHING
           `

			if _, err := db.Exec(query, roomID, date, isBooked, dailyRate); err != nil {
				log.Printf("Error inserting data for room %s on %s: %v", roomID, date.Format("2006-01-02"), err)
				continue
			}
		}

		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM room_bookings WHERE room_id = $1", roomID).Scan(&count); err != nil {
			log.Printf("Error counting records for room %s: %v", roomID, err)
		}
	}

	return roomIDs
}

// round rounds a floating-point number to two decimal places.
// Used for currency calculations in room rates.
//
// Parameters:
//   - num float64: Number to round
//
// Returns:
//   - float64: Number rounded to two decimal places
func round(num float64) float64 {
	return float64(int(num*100)) / 100
}

// main is the entry point of the database setup script.
// It performs the following operations in order:
// 1. Loads environment variables
// 2. Creates database if it doesn't exist
// 3. Establishes database connection
// 4. Creates necessary tables
// 5. Generates and inserts mock data
// 6. Prints generated room IDs
func main() {
	if err := checkAndCreateDatabase(); err != nil {
		log.Fatal(err)
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
			os.Exit(1)
		}
	}()

	if err := checkAndCreateTables(db); err != nil {
		log.Fatal(err)
	}

	roomIDs := generateMockData(db)

	fmt.Println("\nGenerated data with the following room IDs:")
	for _, roomID := range roomIDs {
		fmt.Printf("- %s\n", roomID)
	}
}
