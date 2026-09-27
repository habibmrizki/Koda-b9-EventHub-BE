package services

import (
	"context"
	"errors"
	"strings"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
	"github.com/habibmrizki/BE-EventHub/pkg"
)

type AuthService struct {
	repo       *repositories.AuthRepository
	hashConfig *pkg.HashConfig
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{
		repo:       repo,
		hashConfig: pkg.NewHashConfig(),
	}
}

func (s *AuthService) Register(ctx context.Context, body dto.RegisterRequest) (dto.UserResponse, error) {
	body.FullName = strings.TrimSpace(body.FullName)
	body.Email = strings.TrimSpace(body.Email)

	if err := utils.RegisterValidation(body); err != nil {
		return dto.UserResponse{}, err
	}

	hashedPassword, err := s.hashConfig.GenHash(body.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, body.FullName, body.Email, hashedPassword)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		Location:  user.Location,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, body dto.LoginRequest) (dto.AuthData, error) {
	email := strings.TrimSpace(body.Email)
	password := strings.TrimSpace(body.Password)

	if email == "" || password == "" {
		return dto.AuthData{}, errors.New("email dan password wajib diisi")
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return dto.AuthData{}, errors.New("email atau password salah")
	}

	isMatch, err := s.hashConfig.CompareHashAndPassword(password, user.Password)
	if err != nil || !isMatch {
		return dto.AuthData{}, errors.New("email atau password salah")
	}

	claims := pkg.NewJWTClaims(user.ID, user.Role)
	token, err := claims.GenToken()
	if err != nil {
		return dto.AuthData{}, err
	}

	return dto.AuthData{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			Location:  user.Location,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string, claims *pkg.Claims) error {
	if token == "" {
		return errors.New("token tidak ditemukan")
	}
	if claims == nil {
		return errors.New("claims token tidak valid")
	}

	if claims.ExpiresAt == nil {
		return errors.New("token tidak memiliki masa berlaku")
	}

	utils.AddToBlacklist(token, claims.ExpiresAt.Time)
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, body dto.ForgotPasswordRequest) error {
	email := strings.TrimSpace(body.Email)
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	_, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return errors.New("email tidak terdaftar")
		}
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, body dto.ResetPasswordRequest) error {
	body.Email = strings.TrimSpace(body.Email)

	if err := utils.ResetPasswordValidation(body); err != nil {
		return err
	}

	hashedPassword, err := s.hashConfig.GenHash(body.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePasswordByEmail(ctx, body.Email, hashedPassword)
}
