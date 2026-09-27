package dto

import (
	"time"
)

type CommunityMemberResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar"`
}

type CommunityDiscussionResponse struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	TimeAgo   string    `json:"time_ago,omitempty"`
}

type CommunityResponse struct {
	ID                  int       `json:"id"`
	OrganizerID         int       `json:"organizerId"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	Description         string    `json:"description"`
	CoverImage          string    `json:"cover_image"`
	CoverImageURL       string    `json:"coverImageUrl"`
	Categories          []string  `json:"categories"`
	MembersCount        int       `json:"members_count"`
	UpcomingEventsCount int       `json:"upcoming_events_count"`
	IsJoined            bool      `json:"is_joined"`
	IsActive            bool      `json:"isActive"`
	CreatedAt           time.Time `json:"createdAt"`
}

type CommunityDetailResponse struct {
	CommunityResponse
	Members     []CommunityMemberResponse     `json:"members"`
	Discussions []CommunityDiscussionResponse `json:"discussions"`
}

type CreateDiscussionRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type CommunityFilterQuery struct {
	Search   string `form:"search"`
	Category string `form:"category"`
	Status   string `form:"status"` // 'active', 'all'
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=10"`
}

type CommunityListData struct {
	Communities []CommunityResponse `json:"communities"`
	Total       int                 `json:"total"`
	Page        int                 `json:"page"`
	Limit       int                 `json:"limit"`
}

type PopularCommunityResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	CoverImageURL       string `json:"cover_image_url"`
	TotalMembers        int    `json:"total_members"`
	TotalUpcomingEvents int    `json:"total_upcoming_events"`
	PopularityScore     int    `json:"popularity_score"`
}

type CommunityMemberDetail struct {
	ID            int       `json:"id"`
	FullName      string    `json:"full_name"`
	AvatarURL     string    `json:"avatar_url"`
	Location      string    `json:"location"`
	CommunityRole string    `json:"community_role"`
	JoinedAt      time.Time `json:"joined_at"`
}
