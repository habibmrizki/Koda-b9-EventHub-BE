package dto

import "time"

type TestimonialResponse struct {
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	AuthorTitle string    `json:"author_title"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
	IsFeatured  bool      `json:"is_featured"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTestimonialRequest struct {
	Content     string `json:"content" form:"content" binding:"required"`
	AuthorTitle string `json:"author_title" form:"author_title" binding:"required"`
}
