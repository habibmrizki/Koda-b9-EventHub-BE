	package dto

	import (
		"mime/multipart"
		"time"
	)

	type UserProfileResponse struct {
		ID         int       `json:"id"`
		FullName   string    `json:"fullName"`
		Email      string    `json:"email"`
		Role       string    `json:"role"`
		Location   string    `json:"location"`
		AvatarURL  string    `json:"avatarUrl"`
		Bio        string    `json:"bio"`
		JoinedDate string    `json:"joinedDate"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
	}

	type UpdateProfileRequest struct {
		FullName  *string               `json:"fullName" form:"fullName"`
		Location  *string               `json:"location" form:"location"`
		Bio       *string               `json:"bio" form:"bio"`
		Avatar    *multipart.FileHeader `json:"-" form:"avatar"`
		AvatarURL *string               `json:"avatarUrl" form:"avatarUrl"`
	}

	type ChangePasswordRequest struct {
		OldPassword     string `json:"oldPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
	}
