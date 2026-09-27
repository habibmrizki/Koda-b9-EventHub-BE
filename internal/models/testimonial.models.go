package models

import "time"

type Testimonial struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Content     string    `json:"content"`
	AuthorTitle string    `json:"author_title"`
	IsFeatured  bool      `json:"is_featured"`
	CreatedAt   time.Time `json:"created_at"`
}
