package services

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
)

type CommunityService struct {
	communityRepo *repositories.CommunityRepository
}

func NewCommunityService(communityRepo *repositories.CommunityRepository) *CommunityService {
	return &CommunityService{communityRepo: communityRepo}
}

func (s *CommunityService) GetAllCommunities(ctx context.Context, filter dto.CommunityFilterQuery, currentUserID int) (*dto.CommunityListData, error) {
	communities, total, err := s.communityRepo.GetAll(ctx, filter, currentUserID)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	return &dto.CommunityListData{
		Communities: communities,
		Total:       total,
		Page:        page,
		Limit:       limit,
	}, nil
}

func (s *CommunityService) GetCommunityByID(ctx context.Context, id int, currentUserID int) (*dto.CommunityDetailResponse, error) {
	return s.communityRepo.GetByID(ctx, id, currentUserID)
}

func (s *CommunityService) JoinCommunity(ctx context.Context, communityID int, userID int) error {
	return s.communityRepo.Join(ctx, communityID, userID)
}

func (s *CommunityService) LeaveCommunity(ctx context.Context, communityID int, userID int) error {
	return s.communityRepo.Leave(ctx, communityID, userID)
}

func (s *CommunityService) AddDiscussion(ctx context.Context, communityID int, userID int, content string) error {
	return s.communityRepo.AddDiscussion(ctx, communityID, userID, content)
}

func (s *CommunityService) GetPopularCommunities(ctx context.Context, limit int) ([]dto.PopularCommunityResponse, error) {
	return s.communityRepo.GetPopular(ctx, limit)
}

func (s *CommunityService) GetCommunityMembers(ctx context.Context, communityID int) ([]dto.CommunityMemberDetail, error) {
	return s.communityRepo.GetMembers(ctx, communityID)
}
