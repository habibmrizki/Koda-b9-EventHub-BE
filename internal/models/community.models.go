package models

import "time"

type Community struct {
	ID            int       `json:"id" db:"id"`
	OrganizerID   int       `json:"organizerId" db:"organizer_id"`
	Name          string    `json:"name" db:"name"`
	Slug          string    `json:"slug" db:"slug"`
	Description   string    `json:"description" db:"description"`
	CoverImageURL string    `json:"coverImageUrl" db:"cover_image_url"`
	IsActive      bool      `json:"isActive" db:"is_active"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type CommunityCategory struct {
	CommunityID int `json:"communityId" db:"community_id"`
	CategoryID  int `json:"categoryId" db:"category_id"`
}

type CommunityMember struct {
	ID            int       `json:"id" db:"id"`
	CommunityID   int       `json:"communityId" db:"community_id"`
	UserID        int       `json:"userId" db:"user_id"`
	CommunityRole string    `json:"communityRole" db:"community_role"`
	JoinedAt      time.Time `json:"joinedAt" db:"joined_at"`
}

type CommunityDiscussion struct {
	ID          int       `json:"id" db:"id"`
	CommunityID int       `json:"communityId" db:"community_id"`
	UserID      int       `json:"userId" db:"user_id"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
