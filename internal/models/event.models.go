package models

import "time"

type Event struct {
	ID           int       `json:"id" db:"id"`
	OrganizerID  int       `json:"organizerId" db:"organizer_id"`
	CommunityID  *int      `json:"communityId" db:"community_id"`
	Title        string    `json:"title" db:"title"`
	Slug         string    `json:"slug" db:"slug"`
	Overview     string    `json:"overview" db:"overview"`
	Description  string    `json:"description" db:"description"`
	StartTime    time.Time `json:"startTime" db:"start_time"`
	EndTime      time.Time `json:"endTime" db:"end_time"`
	LocationType string    `json:"locationType" db:"location_type"`
	City         string    `json:"city" db:"city"`
	Address      string    `json:"address" db:"address"`
	Capacity     int       `json:"capacity" db:"capacity"`
	ThumbnailURL string    `json:"thumbnailUrl" db:"thumbnail_url"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type EventRegistration struct {
	ID           int       `json:"id" db:"id"`
	EventID      int       `json:"eventId" db:"event_id"`
	UserID       int       `json:"userId" db:"user_id"`
	Status       string    `json:"status" db:"status"`
	RegisteredAt time.Time `json:"registeredAt" db:"registered_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type EventBookmark struct {
	UserID    int       `json:"userId" db:"user_id"`
	EventID   int       `json:"eventId" db:"event_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}
