package repository

import (
	"airbnb-analytics/internal/database"
	"airbnb-analytics/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// RoomRepository handles database operations for room data
type RoomRepository struct {
	db *sql.DB
}

// NewRoomRepository creates a new repository instance with database connection.
// Returns:
//   - *RoomRepository: New repository instance
func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		db: database.DB,
	}
}

// GetRoomData retrieves room booking data for a given date range.
// Parameters:
//   - roomID string: Room identifier
//   - startDate time.Time: Start of date range
//   - endDate time.Time: End of date range
//
// Returns:
//   - []models.RoomData: Slice of room booking data
//   - error: Any error encountered
func (r *RoomRepository) GetRoomData(roomID string, startDate, endDate time.Time) (roomData []models.RoomData, err error) {
	query := `
        SELECT date::date, is_booked, rate 
        FROM room_bookings 
        WHERE room_id = $1 
        AND date::date >= $2::date 
        AND date::date <= $3::date
        ORDER BY date
    `

	rows, err := r.db.Query(query, roomID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying room data: %v", err)
	}

	// Using named return to handle close error
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing rows: %v", closeErr)
		}
	}()

	var bookings []models.RoomData
	for rows.Next() {
		var booking models.RoomData
		var date time.Time
		if err := rows.Scan(&date, &booking.IsBooked, &booking.Rate); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		booking.Date = date.Format("2006-01-02")
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return bookings, nil
}

// GetAllRoomIDs retrieves all unique room identifiers.
// Returns:
//   - []string: List of room IDs
//   - error: Any error encountered
func (r *RoomRepository) GetAllRoomIDs() (roomIDs []string, err error) {
	query := `SELECT DISTINCT room_id FROM room_bookings ORDER BY room_id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying room IDs: %v", err)
	}

	// Using named return to handle close error
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing rows: %v", closeErr)
		}
	}()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning room ID: %v", err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating room IDs: %v", err)
	}

	return ids, nil
}
