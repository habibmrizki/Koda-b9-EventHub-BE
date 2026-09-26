package services

import (
	"context"
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
