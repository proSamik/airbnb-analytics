package service

import (
	"airbnb-analytics/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type RoomService struct {
	mockAPIURL string
}

func NewRoomService() *RoomService {
	return &RoomService{
		mockAPIURL: "http://localhost:3001/rooms",
	}
}

func (s *RoomService) GetRoomAnalytics(roomID string) (*models.AnalyticsResponse, error) {
	resp, err := http.Get(s.mockAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room data: %v", err)
	}
	defer resp.Body.Close()

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

	// Calculate occupancy and rates separately
	occupancy := calculateMonthlyOccupancy(roomData)
	rateAnalytics := calculateRateAnalytics(roomData)

	return &models.AnalyticsResponse{
		RoomID:           roomID,
		MonthlyOccupancy: occupancy,
		RateAnalytics:    rateAnalytics,
	}, nil
}

// calculateMonthlyOccupancy calculates occupancy for next 5 months
func calculateMonthlyOccupancy(data []models.RoomData) []models.MonthlyOccupancy {
	monthlyStats := make(map[string]struct {
		booked int
		total  int
	})

	currentTime := time.Now()
	fiveMonthsFromNow := currentTime.AddDate(0, 5, 0)

	// Calculate occupancy for each month
	for _, booking := range data {
		bookingDate, err := time.Parse("2006-01-02", booking.Date)
		if err != nil {
			continue
		}

		// Skip if date is before today or after 5 months
		if bookingDate.Before(currentTime) || bookingDate.After(fiveMonthsFromNow) {
			continue
		}

		month := bookingDate.Format("2006-01")
		stats := monthlyStats[month]
		stats.total++
		if booking.IsBooked {
			stats.booked++
		}
		monthlyStats[month] = stats
	}

	// Convert map to sorted slice
	var occupancy []models.MonthlyOccupancy
	for month, stats := range monthlyStats {
		percent := (float64(stats.booked) / float64(stats.total)) * 100
		occupancy = append(occupancy, models.MonthlyOccupancy{
			Month:               month,
			OccupancyPercentage: round(percent),
		})
	}

	// Sort by month
	sort.Slice(occupancy, func(i, j int) bool {
		return occupancy[i].Month < occupancy[j].Month
	})

	return occupancy
}

// calculateRateAnalytics calculates rate statistics for next 30 days
func calculateRateAnalytics(data []models.RoomData) models.RateAnalytics {
	var rates []float64
	currentTime := time.Now()
	thirtyDaysFromNow := currentTime.AddDate(0, 0, 30)

	// Collect rates for next 30 days
	for _, booking := range data {
		bookingDate, err := time.Parse("2006-01-02", booking.Date)
		if err != nil {
			continue
		}

		// Only include rates for next 30 days
		if !bookingDate.Before(currentTime) && !bookingDate.After(thirtyDaysFromNow) {
			rates = append(rates, booking.Rate)
		}
	}

	// Initialize with default values
	rateAnalytics := models.RateAnalytics{
		AverageRate: 0,
		HighestRate: 0,
		LowestRate:  0,
	}

	// Calculate statistics if we have rates
	if len(rates) > 0 {
		var sum float64
		highest := rates[0]
		lowest := rates[0]

		for _, rate := range rates {
			sum += rate
			if rate > highest {
				highest = rate
			}
			if rate < lowest {
				lowest = rate
			}
		}

		rateAnalytics = models.RateAnalytics{
			AverageRate: round(sum / float64(len(rates))),
			HighestRate: round(highest),
			LowestRate:  round(lowest),
		}
	}

	return rateAnalytics
}

// Helper function to round to 2 decimal places
func round(num float64) float64 {
	return float64(int(num*100)) / 100
}
