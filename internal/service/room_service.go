package service

import (
	"airbnb-analytics/internal/models"
	"airbnb-analytics/internal/repository"
	"fmt"
	"time"
)

// RoomService handles room analytics operations and database interactions.
// It processes raw booking data to generate occupancy and rate analytics.
type RoomService struct {
	repo *repository.RoomRepository
}

// NewRoomService creates and returns a new RoomService instance
// with configured repository.
// Returns:
//   - *RoomService: New room service instance configured with repository
func NewRoomService() *RoomService {
	return &RoomService{
		repo: repository.NewRoomRepository(),
	}
}

// GetRoomAnalytics retrieves and processes analytics data for a specific room.
// It calculates occupancy rates for the next 5 months and rate analytics for
// the next 30 days from the current date.
// Parameters:
//   - roomID string: Unique identifier for the room
//
// Returns:
//   - *models.AnalyticsResponse: Processed analytics data containing occupancy and rate statistics
//   - error: Any error encountered during data retrieval or processing
func (s *RoomService) GetRoomAnalytics(roomID string) (response *models.AnalyticsResponse, err error) {
	// Get data for next 5 months from today
	startDate := time.Now().Truncate(24 * time.Hour) // Start from beginning of today
	endDate := startDate.AddDate(0, 5, 0)

	roomData, err := s.repo.GetRoomData(roomID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room data: %v", err)
	}

	if len(roomData) == 0 {
		return nil, fmt.Errorf("room not found")
	}

	occupancy := calculateMonthlyOccupancy(roomData)
	rateAnalytics := calculateRateAnalytics(roomData)

	return &models.AnalyticsResponse{
		RoomID:           roomID,
		MonthlyOccupancy: occupancy,
		RateAnalytics:    rateAnalytics,
	}, nil
}

// GetAllRooms retrieves a list of all available room IDs from the database.
// This can be used to get an overview of all rooms in the system.
// Returns:
//   - []string: List of unique room identifiers
//   - error: Any error encountered during the database operation
func (s *RoomService) GetAllRooms() ([]string, error) {
	roomIDs, err := s.repo.GetAllRoomIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room IDs: %v", err)
	}
	return roomIDs, nil
}
