package models

import "time"

type Users struct {
	ID        int       `json:"id" db:"id"`
	FullName  string    `json:"fullName" db:"full_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	Location  string    `json:"location" db:"location"`
	AvatarURL string    `json:"avatarUrl" db:"avatar_url"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
