package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
)

type EventService struct {
	repo *repositories.EventRepository
}

func NewEventService(repo *repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func parseEventDateTime(date string, timeStr string) (time.Time, error) {
	loc := time.FixedZone("WIB", 7*60*60)
	return time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", date, timeStr), loc)
}

func (s *EventService) formatEventResponse(item repositories.EventDetailQueryResult, isJoined, isBookmarked bool) dto.EventResponse {
	e := item.Event
	firstCategory := ""
	if len(item.Categories) > 0 {
		firstCategory = item.Categories[0]
	}

	locationDisplay := e.City
	if locationDisplay == "" {
		if strings.EqualFold(e.LocationType, "online") {
			locationDisplay = "Online Event"
		} else {
			locationDisplay = e.Address
		}
	}

	resp := dto.EventResponse{
		ID:            e.ID,
		Title:         e.Title,
		Slug:          e.Slug,
		Overview:      e.Overview,
		Description:   e.Description,
		ThumbnailURL:  e.ThumbnailURL,
		CoverImage:    e.ThumbnailURL,
		Categories:    item.Categories,
		Category:      firstCategory,
		LocationType:  e.LocationType,
		City:          e.City,
		Address:       e.Address,
		Location:      locationDisplay,
		Date:          e.StartTime.Format("2006-01-02"),
		StartTime:     e.StartTime.Format("15:04"),
		EndTime:       e.EndTime.Format("15:04"),
		Capacity:      e.Capacity,
		Registered:    item.Registered,
		IsFull:        item.Registered >= e.Capacity,
		IsJoined:      isJoined,
		IsBookmarked:  isBookmarked,
		CommunityName: item.CommunityName,
		Status:        e.Status,
		CreatedAt:     e.CreatedAt,
	}

	if item.Organizer != nil {
		resp.Organizer = &dto.EventOrganizerInfo{
			ID:        item.Organizer.ID,
			FullName:  item.Organizer.FullName,
			Email:     item.Organizer.Email,
			AvatarURL: item.Organizer.AvatarURL,
		}
	}

	return resp
}

func (s *EventService) GetEvents(ctx context.Context, filter dto.EventFilterQuery, userID int) ([]dto.EventResponse, int, error) {
	results, total, err := s.repo.GetEvents(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var formatted []dto.EventResponse
	for _, item := range results {
		isJoined := false
		isBookmarked := false
		if userID > 0 {
			isJoined, _ = s.repo.IsUserJoined(ctx, userID, item.Event.ID)
			isBookmarked, _ = s.repo.IsUserBookmarked(ctx, userID, item.Event.ID)
		}
		formatted = append(formatted, s.formatEventResponse(item, isJoined, isBookmarked))
	}

	return formatted, total, nil
}

func (s *EventService) GetEventDetail(ctx context.Context, eventID int, userID int) (dto.EventResponse, error) {
	item, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return dto.EventResponse{}, err
	}

	isJoined := false
	isBookmarked := false
	if userID > 0 {
		isJoined, _ = s.repo.IsUserJoined(ctx, userID, item.Event.ID)
		isBookmarked, _ = s.repo.IsUserBookmarked(ctx, userID, item.Event.ID)
	}

	return s.formatEventResponse(item, isJoined, isBookmarked), nil
}

func (s *EventService) GetUpcomingEvents(ctx context.Context, limit int, userID int) ([]dto.EventResponse, error) {
	results, err := s.repo.GetUpcomingEvents(ctx, limit)
	if err != nil {
		return nil, err
	}

	var formatted []dto.EventResponse
	for _, item := range results {
		isJoined := false
		isBookmarked := false
		if userID > 0 {
			isJoined, _ = s.repo.IsUserJoined(ctx, userID, item.Event.ID)
			isBookmarked, _ = s.repo.IsUserBookmarked(ctx, userID, item.Event.ID)
		}
		formatted = append(formatted, s.formatEventResponse(item, isJoined, isBookmarked))
	}

	return formatted, nil
}

func (s *EventService) GetMyEvents(ctx context.Context, userID int) ([]dto.EventResponse, error) {
	results, err := s.repo.GetMyEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	var formatted []dto.EventResponse
	for _, item := range results {
		isBookmarked, _ := s.repo.IsUserBookmarked(ctx, userID, item.Event.ID)
		formatted = append(formatted, s.formatEventResponse(item, true, isBookmarked))
	}

	return formatted, nil
}

func (s *EventService) JoinEvent(ctx context.Context, userID, eventID int) error {
	return s.repo.JoinEvent(ctx, userID, eventID)
}

func (s *EventService) LeaveEvent(ctx context.Context, userID, eventID int) error {
	return s.repo.LeaveEvent(ctx, userID, eventID)
}

func (s *EventService) CreateEvent(ctx context.Context, organizerID int, req dto.CreateEventRequest, thumbnailURL string) (dto.EventResponse, error) {
	slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
	slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())

	startTime, err := parseEventDateTime(req.EventDate, req.StartTime)
	if err != nil {
		startTime = time.Now().In(time.FixedZone("WIB", 7*60*60)).Add(24 * time.Hour)
	}

	endTime, err := parseEventDateTime(req.EventDate, req.EndTime)
	if err != nil {
		endTime = startTime.Add(2 * time.Hour)
	}

	if endTime.Before(startTime) {
		endTime = startTime.Add(2 * time.Hour)
	}

	eventID, err := s.repo.CreateEvent(ctx, organizerID, req, slug, startTime, endTime, thumbnailURL)
	if err != nil {
		return dto.EventResponse{}, err
	}

	return s.GetEventDetail(ctx, eventID, organizerID)
}

func (s *EventService) UpdateEvent(ctx context.Context, eventID int, organizerID int, req dto.UpdateEventRequest, thumbnailURL *string) (dto.EventResponse, error) {
	var startTime *time.Time
	var endTime *time.Time

	if req.EventDate != nil && req.StartTime != nil {
		if st, err := parseEventDateTime(*req.EventDate, *req.StartTime); err == nil {
			startTime = &st
		}
	}

	if req.EventDate != nil && req.EndTime != nil {
		if et, err := parseEventDateTime(*req.EventDate, *req.EndTime); err == nil {
			endTime = &et
		}
	}

	if startTime != nil && endTime != nil && endTime.Before(*startTime) {
		return dto.EventResponse{}, fmt.Errorf("endTime tidak boleh lebih kecil dari startTime")
	}

	if err := s.repo.UpdateEvent(ctx, eventID, organizerID, req, startTime, endTime, thumbnailURL); err != nil {
		return dto.EventResponse{}, err
	}

	return s.GetEventDetail(ctx, eventID, organizerID)
}
