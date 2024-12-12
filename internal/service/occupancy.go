package service

import (
	"airbnb-analytics/internal/models"
	"sort"
	"time"
)

// calculateMonthlyOccupancy processes room data to calculate occupancy rates
// for the next 5 months.
// Parameters:
//   - data []models.RoomData: Slice of room booking data
//
// Returns:
//   - []models.MonthlyOccupancy: Slice of monthly occupancy statistics
func calculateMonthlyOccupancy(data []models.RoomData) []models.MonthlyOccupancy {
	monthlyStats := make(map[string]struct {
		booked int
		total  int
	})

	currentTime := time.Now()
	fiveMonthsFromNow := currentTime.AddDate(0, 5, 0)

	// Calculate monthly statistics
	for _, booking := range data {
		bookingDate, err := time.Parse("2006-01-02", booking.Date)
		if err != nil {
			continue
		}

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
