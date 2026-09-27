package services

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetMyNotifications(ctx context.Context, userID int) ([]dto.NotificationResponse, error) {
	return s.repo.GetMyNotifications(ctx, userID)
}
