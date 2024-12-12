package service

import (
	"airbnb-analytics/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// RoomService handles room analytics operations and API interactions.
type RoomService struct {
	mockAPIURL string
}

// NewRoomService creates and returns a new RoomService instance
// with configured mock API URL.
// Returns:
//   - *RoomService: New room service instance
func NewRoomService() *RoomService {
	return &RoomService{
		mockAPIURL: "http://localhost:3001/rooms",
	}
}

// GetRoomAnalytics retrieves and processes analytics data for a specific room.
// Parameters:
//   - roomID string: Unique identifier for the room
//
// Returns:
//   - *models.AnalyticsResponse: Processed analytics data
//   - error: Any error encountered during the operation
func (s *RoomService) GetRoomAnalytics(roomID string) (response *models.AnalyticsResponse, err error) {
	resp, err := http.Get(s.mockAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room data: %v", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var roomsData map[string][]models.RoomData
	if err := json.NewDecoder(resp.Body).Decode(&roomsData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	roomData, exists := roomsData[roomID]
	if !exists {
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
