package services

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
)

type DashboardService struct {
	repo *repositories.DashboardRepository
}

func NewDashboardService(repo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetAdminDashboard(ctx context.Context) (dto.AdminDashboardResponse, error) {
	return s.repo.GetAdminDashboard(ctx)
}
