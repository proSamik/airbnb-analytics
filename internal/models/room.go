package models

type RoomData struct {
	Date     string  `json:"date"`
	IsBooked bool    `json:"is_booked"`
	Rate     float64 `json:"rate"`
}

type AnalyticsResponse struct {
	RoomID           string             `json:"room_id"`
	MonthlyOccupancy []MonthlyOccupancy `json:"monthly_occupancy"`
	RateAnalytics    RateAnalytics      `json:"rate_analytics"`
}

type MonthlyOccupancy struct {
	Month               string  `json:"month"`
	OccupancyPercentage float64 `json:"occupancy_percentage"`
}

type RateAnalytics struct {
	AverageRate float64 `json:"average_rate"`
	HighestRate float64 `json:"highest_rate"`
	LowestRate  float64 `json:"lowest_rate"`
}
