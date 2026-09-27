package services

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
)

type TestimonialService struct {
	repo *repositories.TestimonialRepository
}

func NewTestimonialService(repo *repositories.TestimonialRepository) *TestimonialService {
	return &TestimonialService{repo: repo}
}

func (s *TestimonialService) GetFeaturedTestimonials(ctx context.Context, limit int) ([]dto.TestimonialResponse, error) {
	return s.repo.GetFeatured(ctx, limit)
}

func (s *TestimonialService) GetAllTestimonials(ctx context.Context) ([]dto.TestimonialResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *TestimonialService) CreateTestimonial(ctx context.Context, userID int, req dto.CreateTestimonialRequest) (dto.TestimonialResponse, error) {
	return s.repo.Create(ctx, userID, req)
}
