package models

// RoomData represents the booking information for a single day of a room.
// It contains date, booking status and rate information.
type RoomData struct {
	// Date represents the booking date in "YYYY-MM-DD" format
	Date string `json:"date"`
	// IsBooked indicates whether the room is booked for this date
	IsBooked bool `json:"is_booked"`
	// Rate represents the room rate for this date in the local currency
	Rate float64 `json:"rate"`
}

// AnalyticsResponse represents the complete analytics response for a room.
// It includes the room identifier, occupancy data and rate analytics.
type AnalyticsResponse struct {
	// RoomID uniquely identifies the room
	RoomID string `json:"room_id"`
	// MonthlyOccupancy contains occupancy data for upcoming months
	MonthlyOccupancy []MonthlyOccupancy `json:"monthly_occupancy"`
	// RateAnalytics contains statistical analysis of room rates
	RateAnalytics RateAnalytics `json:"rate_analytics"`
}

// MonthlyOccupancy represents the occupancy statistics for a single month.
type MonthlyOccupancy struct {
	// Month represents the month in "YYYY-MM" format
	Month string `json:"month"`
	// OccupancyPercentage represents the percentage of days booked in this month
	OccupancyPercentage float64 `json:"occupancy_percentage"`
}

// RateAnalytics represents statistical analysis of room rates.
// All rates are in the local currency.
type RateAnalytics struct {
	// AverageRate represents the mean rate across the analyzed period
	AverageRate float64 `json:"average_rate"`
	// HighestRate represents the maximum rate in the analyzed period
	HighestRate float64 `json:"highest_rate"`
	// LowestRate represents the minimum rate in the analyzed period
	LowestRate float64 `json:"lowest_rate"`
}
