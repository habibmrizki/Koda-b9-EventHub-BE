package dto

type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Location  string `json:"location"`
	AvatarURL string `json:"avatarUrl"`
	Bio       string `json:"bio"`
}

type AuthData struct {
	Token string       `json:"token,omitempty"`
	User  UserResponse `json:"user"`
}
