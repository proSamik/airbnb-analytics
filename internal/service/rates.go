package service

import (
	"airbnb-analytics/internal/models"
	"time"
)

// calculateRateAnalytics processes room data to calculate rate statistics
// for the next 30 days.
// Parameters:
//   - data []models.RoomData: Slice of room booking data
//
// Returns:
//   - models.RateAnalytics: Calculated rate statistics
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

		if !bookingDate.Before(currentTime) && !bookingDate.After(thirtyDaysFromNow) {
			rates = append(rates, booking.Rate)
		}
	}

	rateAnalytics := models.RateAnalytics{
		AverageRate: 0,
		HighestRate: 0,
		LowestRate:  0,
	}

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

// round rounds a float64 to 2 decimal places.
// Parameters:
//   - num float64: Number to round
//
// Returns:
//   - float64: Rounded number
func round(num float64) float64 {
	return float64(int(num*100)) / 100
}
