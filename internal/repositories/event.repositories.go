package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEventNotFound  = errors.New("event not found")
	ErrEventFull      = errors.New("event is full, capacity reached")
	ErrAlreadyJoined  = errors.New("already joined this event")
	ErrNotJoinedEvent = errors.New("you have not joined this event")
	ErrUnauthorized   = errors.New("unauthorized to modify this event")
)

type EventDetailQueryResult struct {
	Event         models.Event
	CommunityName string
	Organizer     *models.Users
	Registered    int
	Categories    []string
}

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetEvents(ctx context.Context, filter dto.EventFilterQuery) ([]EventDetailQueryResult, int, error) {
	conditions := []string{"e.status = 'published'"}
	args := []interface{}{}
	argPos := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(LOWER(e.title) LIKE $%d OR LOWER(e.description) LIKE $%d OR LOWER(e.overview) LIKE $%d)", argPos, argPos, argPos))
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		argPos++
	}

	if filter.Location != "" {
		conditions = append(conditions, fmt.Sprintf("(LOWER(e.city) LIKE $%d OR LOWER(e.address) LIKE $%d)", argPos, argPos))
		args = append(args, "%"+strings.ToLower(filter.Location)+"%")
		argPos++
	}

	if filter.LocationType != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.location_type) = $%d", argPos))
		args = append(args, strings.ToLower(filter.LocationType))
		argPos++
	}

	categoryFilter := ""
	if filter.Category != "" {
		categoryFilter = fmt.Sprintf("HAVING $%d = ANY(ARRAY_AGG(LOWER(cat.name)))", argPos)
		args = append(args, strings.ToLower(filter.Category))
		argPos++
	}

	whereClause := strings.Join(conditions, " AND ")

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT 
			e.id, e.organizer_id, e.community_id, e.title, e.slug, COALESCE(e.overview, ''), e.description,
			e.start_time, e.end_time, e.location_type, COALESCE(e.city, ''), COALESCE(e.address, ''),
			e.capacity, COALESCE(e.thumbnail_url, ''), e.status, e.created_at, e.updated_at,
			COALESCE(c.name, '') as community_name,
			COALESCE(u.id, 0), COALESCE(u.full_name, ''), COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
			(SELECT COUNT(*) FROM event_registrations er WHERE er.event_id = e.id AND er.status = 'registered') as registered_count,
			COALESCE(ARRAY_AGG(cat.name) FILTER (WHERE cat.name IS NOT NULL), '{}') as categories
		FROM events e
		LEFT JOIN users u ON e.organizer_id = u.id
		LEFT JOIN communities c ON e.community_id = c.id
		LEFT JOIN event_categories ec ON e.id = ec.event_id
		LEFT JOIN categories cat ON ec.category_id = cat.id
		WHERE %s
		GROUP BY e.id, c.name, u.id, u.full_name, u.email, u.avatar_url
		%s
		ORDER BY e.start_time ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, categoryFilter, argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []EventDetailQueryResult
	for rows.Next() {
		var item EventDetailQueryResult
		var org models.Users
		var cats []string
		if err := rows.Scan(
			&item.Event.ID,
			&item.Event.OrganizerID,
			&item.Event.CommunityID,
			&item.Event.Title,
			&item.Event.Slug,
			&item.Event.Overview,
			&item.Event.Description,
			&item.Event.StartTime,
			&item.Event.EndTime,
			&item.Event.LocationType,
			&item.Event.City,
			&item.Event.Address,
			&item.Event.Capacity,
			&item.Event.ThumbnailURL,
			&item.Event.Status,
			&item.Event.CreatedAt,
			&item.Event.UpdatedAt,
			&item.CommunityName,
			&org.ID,
			&org.FullName,
			&org.Email,
			&org.AvatarURL,
			&item.Registered,
			&cats,
		); err != nil {
			return nil, 0, err
		}
		if org.ID != 0 {
			item.Organizer = &org
		}
		item.Categories = cats
		results = append(results, item)
	}

	totalQuery := fmt.Sprintf("SELECT COUNT(DISTINCT e.id) FROM events e WHERE %s", whereClause)
	var total int
	_ = r.db.QueryRow(ctx, totalQuery, args[:len(conditions)-1]...).Scan(&total)
	if total == 0 {
		total = len(results)
	}

	return results, total, nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, id int) (EventDetailQueryResult, error) {
	query := `
		SELECT 
			e.id, e.organizer_id, e.community_id, e.title, e.slug, COALESCE(e.overview, ''), e.description,
			e.start_time, e.end_time, e.location_type, COALESCE(e.city, ''), COALESCE(e.address, ''),
			e.capacity, COALESCE(e.thumbnail_url, ''), e.status, e.created_at, e.updated_at,
			COALESCE(c.name, '') as community_name,
			COALESCE(u.id, 0), COALESCE(u.full_name, ''), COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
			(SELECT COUNT(*) FROM event_registrations er WHERE er.event_id = e.id AND er.status = 'registered') as registered_count,
			COALESCE(ARRAY_AGG(cat.name) FILTER (WHERE cat.name IS NOT NULL), '{}') as categories
		FROM events e
		LEFT JOIN users u ON e.organizer_id = u.id
		LEFT JOIN communities c ON e.community_id = c.id
		LEFT JOIN event_categories ec ON e.id = ec.event_id
		LEFT JOIN categories cat ON ec.category_id = cat.id
		WHERE e.id = $1
		GROUP BY e.id, c.name, u.id, u.full_name, u.email, u.avatar_url
		LIMIT 1
	`

	var item EventDetailQueryResult
	var org models.Users
	var cats []string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.Event.ID,
		&item.Event.OrganizerID,
		&item.Event.CommunityID,
		&item.Event.Title,
		&item.Event.Slug,
		&item.Event.Overview,
		&item.Event.Description,
		&item.Event.StartTime,
		&item.Event.EndTime,
		&item.Event.LocationType,
		&item.Event.City,
		&item.Event.Address,
		&item.Event.Capacity,
		&item.Event.ThumbnailURL,
		&item.Event.Status,
		&item.Event.CreatedAt,
		&item.Event.UpdatedAt,
		&item.CommunityName,
		&org.ID,
		&org.FullName,
		&org.Email,
		&org.AvatarURL,
		&item.Registered,
		&cats,
	)
	if err != nil {
		return EventDetailQueryResult{}, ErrEventNotFound
	}

	if org.ID != 0 {
		item.Organizer = &org
	}
	item.Categories = cats

	return item, nil
}

func (r *EventRepository) GetUpcomingEvents(ctx context.Context, limit int) ([]EventDetailQueryResult, error) {
	if limit <= 0 {
		limit = 6
	}

	query := `
		SELECT 
			e.id, e.organizer_id, e.community_id, e.title, e.slug, COALESCE(e.overview, ''), e.description,
			e.start_time, e.end_time, e.location_type, COALESCE(e.city, ''), COALESCE(e.address, ''),
			e.capacity, COALESCE(e.thumbnail_url, ''), e.status, e.created_at, e.updated_at,
			COALESCE(c.name, '') as community_name,
			COALESCE(u.id, 0), COALESCE(u.full_name, ''), COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
			(SELECT COUNT(*) FROM event_registrations er WHERE er.event_id = e.id AND er.status = 'registered') as registered_count,
			COALESCE(ARRAY_AGG(cat.name) FILTER (WHERE cat.name IS NOT NULL), '{}') as categories
		FROM events e
		LEFT JOIN users u ON e.organizer_id = u.id
		LEFT JOIN communities c ON e.community_id = c.id
		LEFT JOIN event_categories ec ON e.id = ec.event_id
		LEFT JOIN categories cat ON ec.category_id = cat.id
		WHERE e.status = 'published' AND e.start_time >= NOW()
		GROUP BY e.id, c.name, u.id, u.full_name, u.email, u.avatar_url
		ORDER BY e.start_time ASC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []EventDetailQueryResult
	for rows.Next() {
		var item EventDetailQueryResult
		var org models.Users
		var cats []string
		if err := rows.Scan(
			&item.Event.ID,
			&item.Event.OrganizerID,
			&item.Event.CommunityID,
			&item.Event.Title,
			&item.Event.Slug,
			&item.Event.Overview,
			&item.Event.Description,
			&item.Event.StartTime,
			&item.Event.EndTime,
			&item.Event.LocationType,
			&item.Event.City,
			&item.Event.Address,
			&item.Event.Capacity,
			&item.Event.ThumbnailURL,
			&item.Event.Status,
			&item.Event.CreatedAt,
			&item.Event.UpdatedAt,
			&item.CommunityName,
			&org.ID,
			&org.FullName,
			&org.Email,
			&org.AvatarURL,
			&item.Registered,
			&cats,
		); err != nil {
			return nil, err
		}
		if org.ID != 0 {
			item.Organizer = &org
		}
		item.Categories = cats
		results = append(results, item)
	}

	return results, nil
}

func (r *EventRepository) GetMyEvents(ctx context.Context, userID int) ([]EventDetailQueryResult, error) {
	query := `
		SELECT 
			e.id, e.organizer_id, e.community_id, e.title, e.slug, COALESCE(e.overview, ''), e.description,
			e.start_time, e.end_time, e.location_type, COALESCE(e.city, ''), COALESCE(e.address, ''),
			e.capacity, COALESCE(e.thumbnail_url, ''), e.status, e.created_at, e.updated_at,
			COALESCE(c.name, '') as community_name,
			COALESCE(u.id, 0), COALESCE(u.full_name, ''), COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
			(SELECT COUNT(*) FROM event_registrations er WHERE er.event_id = e.id AND er.status = 'registered') as registered_count,
			COALESCE(ARRAY_AGG(cat.name) FILTER (WHERE cat.name IS NOT NULL), '{}') as categories
		FROM events e
		INNER JOIN event_registrations er ON e.id = er.event_id AND er.user_id = $1 AND er.status = 'registered'
		LEFT JOIN users u ON e.organizer_id = u.id
		LEFT JOIN communities c ON e.community_id = c.id
		LEFT JOIN event_categories ec ON e.id = ec.event_id
		LEFT JOIN categories cat ON ec.category_id = cat.id
		GROUP BY e.id, c.name, u.id, u.full_name, u.email, u.avatar_url
		ORDER BY e.start_time ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []EventDetailQueryResult
	for rows.Next() {
		var item EventDetailQueryResult
		var org models.Users
		var cats []string
		if err := rows.Scan(
			&item.Event.ID,
			&item.Event.OrganizerID,
			&item.Event.CommunityID,
			&item.Event.Title,
			&item.Event.Slug,
			&item.Event.Overview,
			&item.Event.Description,
			&item.Event.StartTime,
			&item.Event.EndTime,
			&item.Event.LocationType,
			&item.Event.City,
			&item.Event.Address,
			&item.Event.Capacity,
			&item.Event.ThumbnailURL,
			&item.Event.Status,
			&item.Event.CreatedAt,
			&item.Event.UpdatedAt,
			&item.CommunityName,
			&org.ID,
			&org.FullName,
			&org.Email,
			&org.AvatarURL,
			&item.Registered,
			&cats,
		); err != nil {
			return nil, err
		}
		if org.ID != 0 {
			item.Organizer = &org
		}
		item.Categories = cats
		results = append(results, item)
	}

	return results, nil
}

func (r *EventRepository) IsUserJoined(ctx context.Context, userID, eventID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM event_registrations WHERE user_id = $1 AND event_id = $2 AND status = 'registered')"
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, eventID).Scan(&exists)
	return exists, err
}

func (r *EventRepository) IsUserBookmarked(ctx context.Context, userID, eventID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM event_bookmarks WHERE user_id = $1 AND event_id = $2)"
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, eventID).Scan(&exists)
	return exists, err
}

func (r *EventRepository) JoinEvent(ctx context.Context, userID, eventID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var capacity int
	err = tx.QueryRow(ctx, "SELECT capacity FROM events WHERE id = $1 FOR UPDATE", eventID).Scan(&capacity)
	if err != nil {
		return ErrEventNotFound
	}

	var registeredCount int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM event_registrations WHERE event_id = $1 AND status = 'registered'", eventID).Scan(&registeredCount)
	if err != nil {
		return err
	}

	if registeredCount >= capacity {
		return ErrEventFull
	}

	var exists bool
	_ = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM event_registrations WHERE user_id = $1 AND event_id = $2 AND status = 'registered')", userID, eventID).Scan(&exists)
	if exists {
		return ErrAlreadyJoined
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO event_registrations (event_id, user_id, status, registered_at, updated_at) 
		VALUES ($1, $2, 'registered', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, eventID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *EventRepository) LeaveEvent(ctx context.Context, userID, eventID int) error {
	res, err := r.db.Exec(ctx, "DELETE FROM event_registrations WHERE event_id = $1 AND user_id = $2", eventID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotJoinedEvent
	}

	return nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, organizerID int, req dto.CreateEventRequest, slug string, startTime, endTime time.Time, thumbnailURL string) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	locType := req.LocationType
	if strings.EqualFold(locType, "In Person") {
		locType = "offline"
	} else if strings.EqualFold(locType, "Online") {
		locType = "online"
	}

	query := `
		INSERT INTO events (
			organizer_id, community_id, title, slug, overview, description,
			start_time, end_time, location_type, city, address, capacity,
			thumbnail_url, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 'published', NOW(), NOW())
		RETURNING id
	`

	var eventID int
	err = tx.QueryRow(ctx, query,
		organizerID,
		req.CommunityID,
		req.Title,
		slug,
		req.Overview,
		req.Description,
		startTime,
		endTime,
		locType,
		req.City,
		req.Address,
		req.Capacity,
		thumbnailURL,
	).Scan(&eventID)
	if err != nil {
		return 0, err
	}

	// Link categories
	for _, catName := range req.Categories {
		catName = strings.TrimSpace(catName)
		if catName == "" {
			continue
		}
		catSlug := strings.ToLower(strings.ReplaceAll(catName, " ", "-"))
		var catID int
		err = tx.QueryRow(ctx, `
			INSERT INTO categories (name, slug)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, catName, catSlug).Scan(&catID)
		if err == nil && catID > 0 {
			_, _ = tx.Exec(ctx, `
				INSERT INTO event_categories (event_id, category_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, eventID, catID)
		}
	}

	return eventID, tx.Commit(ctx)
}

func (r *EventRepository) UpdateEvent(ctx context.Context, eventID int, organizerID int, req dto.UpdateEventRequest, startTime, endTime *time.Time, thumbnailURL *string) error {
	var currentOrganizerID int
	err := r.db.QueryRow(ctx, "SELECT organizer_id FROM events WHERE id = $1", eventID).Scan(&currentOrganizerID)
	if err != nil {
		return ErrEventNotFound
	}

	if currentOrganizerID != organizerID {
		return ErrUnauthorized
	}

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Title))
		argPos++
	}
	if req.Overview != nil {
		setClauses = append(setClauses, fmt.Sprintf("overview = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Overview))
		argPos++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Description))
		argPos++
	}
	if startTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_time = $%d", argPos))
		args = append(args, *startTime)
		argPos++
	}
	if endTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_time = $%d", argPos))
		args = append(args, *endTime)
		argPos++
	}
	if req.LocationType != nil {
		locType := *req.LocationType
		if strings.EqualFold(locType, "In Person") {
			locType = "offline"
		}
		setClauses = append(setClauses, fmt.Sprintf("location_type = $%d", argPos))
		args = append(args, locType)
		argPos++
	}
	if req.City != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.City))
		argPos++
	}
	if req.Address != nil {
		setClauses = append(setClauses, fmt.Sprintf("address = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Address))
		argPos++
	}
	if req.Capacity != nil {
		setClauses = append(setClauses, fmt.Sprintf("capacity = $%d", argPos))
		args = append(args, *req.Capacity)
		argPos++
	}
	if thumbnailURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("thumbnail_url = $%d", argPos))
		args = append(args, *thumbnailURL)
		argPos++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Status))
		argPos++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`UPDATE events SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), argPos)
	args = append(args, eventID)

	_, err = r.db.Exec(ctx, query, args...)
	return err
}
