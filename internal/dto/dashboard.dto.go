package dto

import "time"

// Organizer Dashboard DTOs
type OrganizerMetric struct {
	TotalEvents    int     `json:"total_events"`
	TotalAttendees int     `json:"total_attendees"`
	AvgFillRate    float64 `json:"avg_fill_rate"`
}

type OrganizerEventHistory struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	StartTime       time.Time `json:"start_time"`
	LocationType    string    `json:"location_type"`
	City            string    `json:"city"`
	Capacity        int       `json:"capacity"`
	Status          string    `json:"status"`
	TotalRegistered int       `json:"total_registered"`
}

type OrganizerDashboardResponse struct {
	Metrics OrganizerMetric         `json:"metrics"`
	Events  []OrganizerEventHistory `json:"events"`
}

// Admin Dashboard DTO
type AdminDashboardResponse struct {
	TotalUsers         int     `json:"total_users"`
	TotalEvents        int     `json:"total_events"`
	ActiveEvents       int     `json:"active_events"`
	TotalCommunities   int     `json:"total_communities"`
	GlobalAvgFillRate  float64 `json:"global_avg_fill_rate"`
}
