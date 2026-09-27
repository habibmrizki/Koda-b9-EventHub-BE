package dto

import (
	"mime/multipart"
	"time"
)

type EventOrganizerInfo struct {
	ID        int    `json:"id"`
	FullName  string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type EventResponse struct {
	ID            int                 `json:"id"`
	Title         string              `json:"title"`
	Slug          string              `json:"slug"`
	Overview      string              `json:"overview,omitempty"`
	Description   string              `json:"description"`
	ThumbnailURL  string              `json:"thumbnailUrl"`
	CoverImage    string              `json:"coverImage"`
	Categories    []string            `json:"categories"`
	Category      string              `json:"category,omitempty"`
	LocationType  string              `json:"locationType"` // 'online' / 'offline'
	City          string              `json:"city"`
	Address       string              `json:"address,omitempty"`
	Location      string              `json:"location"`
	Date          string              `json:"date"`
	StartTime     string              `json:"startTime"`
	EndTime       string              `json:"endTime"`
	Capacity      int                 `json:"capacity"`
	Registered    int                 `json:"registered"`
	IsFull        bool                `json:"isFull"`
	IsJoined      bool                `json:"isJoined"`
	IsBookmarked  bool                `json:"isBookmarked"`
	CommunityName string              `json:"communityName,omitempty"`
	Organizer     *EventOrganizerInfo `json:"organizer,omitempty"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"createdAt"`
}

type EventListData struct {
	Events []EventResponse `json:"events"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type CreateEventRequest struct {
	Title        string                `json:"title" form:"title" binding:"required"`
	Overview     string                `json:"overview" form:"overview"`
	Description  string                `json:"description" form:"description" binding:"required"`
	EventDate    string                `json:"eventDate" form:"eventDate" binding:"required"`       // e.g. "2026-10-20"
	StartTime    string                `json:"startTime" form:"startTime" binding:"required"`       // e.g. "09:00"
	EndTime      string                `json:"endTime" form:"endTime" binding:"required"`           // e.g. "17:00"
	LocationType string                `json:"locationType" form:"locationType" binding:"required"` // 'online' / 'offline' / 'In Person'
	City         string                `json:"city" form:"city"`
	Address      string                `json:"address" form:"address"`
	Capacity     int                   `json:"capacity" form:"capacity" binding:"required,min=1"`
	CommunityID  *int                  `json:"communityId" form:"communityId"`
	Categories   []string              `json:"categories" form:"categories"`
	Thumbnail    *multipart.FileHeader `json:"-" form:"thumbnail"`
}

type UpdateEventRequest struct {
	Title        *string               `json:"title" form:"title"`
	Overview     *string               `json:"overview" form:"overview"`
	Description  *string               `json:"description" form:"description"`
	EventDate    *string               `json:"eventDate" form:"eventDate"`
	StartTime    *string               `json:"startTime" form:"startTime"`
	EndTime      *string               `json:"endTime" form:"endTime"`
	LocationType *string               `json:"locationType" form:"locationType"`
	City         *string               `json:"city" form:"city"`
	Address      *string               `json:"address" form:"address"`
	Capacity     *int                  `json:"capacity" form:"capacity"`
	Status       *string               `json:"status" form:"status"`
	CommunityID  *int                  `json:"communityId" form:"communityId"`
	Thumbnail    *multipart.FileHeader `json:"-" form:"thumbnail"`
}

type UploadThumbnailResponse struct {
	ThumbnailURL string `json:"thumbnailUrl"`
}

type EventFilterQuery struct {
	Search       string `form:"search"`
	Category     string `form:"category"`
	Location     string `form:"location"`
	LocationType string `form:"locationType"` // 'online' / 'offline'
	Page         int    `form:"page,default=1"`
	Limit        int    `form:"limit,default=10"`
}
