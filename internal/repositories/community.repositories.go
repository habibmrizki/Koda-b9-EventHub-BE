package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepository struct {
	db *pgxpool.Pool
}

func NewCommunityRepository(db *pgxpool.Pool) *CommunityRepository {
	return &CommunityRepository{db: db}
}

func (r *CommunityRepository) GetAll(ctx context.Context, filter dto.CommunityFilterQuery, currentUserID int) ([]dto.CommunityResponse, int, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(c.name ILIKE $%d OR c.description ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	if filter.Status != "" && strings.ToLower(filter.Status) != "all" {
		if strings.ToLower(filter.Status) == "active" {
			conditions = append(conditions, fmt.Sprintf("c.is_active = $%d", argPos))
			args = append(args, true)
			argPos++
		}
	}

	if filter.Category != "" && strings.ToLower(filter.Category) != "all" && strings.ToLower(filter.Category) != "all categories" {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM community_categories cc
			JOIN categories cat ON cc.category_id = cat.id
			WHERE cc.community_id = c.id AND (cat.name ILIKE $%d OR cat.slug ILIKE $%d)
		)`, argPos, argPos))
		args = append(args, "%"+filter.Category+"%")
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM communities c %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	// Query communities
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.organizer_id, c.name, c.slug, COALESCE(c.description, ''), 
			COALESCE(c.cover_image_url, ''), c.is_active, c.created_at,
			(SELECT COUNT(*) FROM community_members cm WHERE cm.community_id = c.id) as members_count,
			(SELECT COUNT(*) FROM events e WHERE e.community_id = c.id AND e.start_time > NOW()) as upcoming_events_count,
			EXISTS(SELECT 1 FROM community_members cm WHERE cm.community_id = c.id AND cm.user_id = $%d) as is_joined
		FROM communities c
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, argPos, whereClause, argPos+1, argPos+2)

	queryArgs := append(args, currentUserID, filter.Limit, offset)
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var communities []dto.CommunityResponse
	for rows.Next() {
		var comm dto.CommunityResponse
		err := rows.Scan(
			&comm.ID, &comm.OrganizerID, &comm.Name, &comm.Slug, &comm.Description,
			&comm.CoverImageURL, &comm.IsActive, &comm.CreatedAt,
			&comm.MembersCount, &comm.UpcomingEventsCount, &comm.IsJoined,
		)
		if err != nil {
			return nil, 0, err
		}
		comm.CoverImage = comm.CoverImageURL
		communities = append(communities, comm)
	}

	// Get categories for each community
	for i := range communities {
		catQuery := `
			SELECT cat.name 
			FROM categories cat
			JOIN community_categories cc ON cat.id = cc.category_id
			WHERE cc.community_id = $1
		`
		catRows, err := r.db.Query(ctx, catQuery, communities[i].ID)
		if err == nil {
			var cats []string
			for catRows.Next() {
				var catName string
				if err := catRows.Scan(&catName); err == nil {
					cats = append(cats, catName)
				}
			}
			catRows.Close()
			communities[i].Categories = cats
		}
	}

	return communities, total, nil
}

func (r *CommunityRepository) GetByID(ctx context.Context, id int, currentUserID int) (*dto.CommunityDetailResponse, error) {
	query := `
		SELECT 
			c.id, c.organizer_id, c.name, c.slug, COALESCE(c.description, ''), 
			COALESCE(c.cover_image_url, ''), c.is_active, c.created_at,
			(SELECT COUNT(*) FROM community_members cm WHERE cm.community_id = c.id) as members_count,
			(SELECT COUNT(*) FROM events e WHERE e.community_id = c.id AND e.start_time > NOW()) as upcoming_events_count,
			EXISTS(SELECT 1 FROM community_members cm WHERE cm.community_id = c.id AND cm.user_id = $2) as is_joined
		FROM communities c
		WHERE c.id = $1
	`

	var detail dto.CommunityDetailResponse
	err := r.db.QueryRow(ctx, query, id, currentUserID).Scan(
		&detail.ID, &detail.OrganizerID, &detail.Name, &detail.Slug, &detail.Description,
		&detail.CoverImageURL, &detail.IsActive, &detail.CreatedAt,
		&detail.MembersCount, &detail.UpcomingEventsCount, &detail.IsJoined,
	)
	if err != nil {
		return nil, err
	}
	detail.CoverImage = detail.CoverImageURL

	// Get categories
	catQuery := `
		SELECT cat.name 
		FROM categories cat
		JOIN community_categories cc ON cat.id = cc.category_id
		WHERE cc.community_id = $1
	`
	catRows, err := r.db.Query(ctx, catQuery, id)
	if err == nil {
		var cats []string
		for catRows.Next() {
			var catName string
			if err := catRows.Scan(&catName); err == nil {
				cats = append(cats, catName)
			}
		}
		catRows.Close()
		detail.Categories = cats
	}

	// Get members
	memQuery := `
		SELECT cm.id, cm.user_id, u.full_name, cm.community_role, COALESCE(u.avatar_url, '')
		FROM community_members cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.community_id = $1
		ORDER BY cm.joined_at ASC
		LIMIT 20
	`
	memRows, err := r.db.Query(ctx, memQuery, id)
	if err == nil {
		var members []dto.CommunityMemberResponse
		for memRows.Next() {
			var m dto.CommunityMemberResponse
			if err := memRows.Scan(&m.ID, &m.UserID, &m.Name, &m.Role, &m.AvatarURL); err == nil {
				members = append(members, m)
			}
		}
		memRows.Close()
		detail.Members = members
	}

	// Get discussions
	discQuery := `
		SELECT cd.id, u.full_name, COALESCE(u.avatar_url, ''), cd.content, cd.created_at
		FROM community_discussions cd
		JOIN users u ON cd.user_id = u.id
		WHERE cd.community_id = $1
		ORDER BY cd.created_at DESC
		LIMIT 50
	`
	discRows, err := r.db.Query(ctx, discQuery, id)
	if err == nil {
		var discussions []dto.CommunityDiscussionResponse
		for discRows.Next() {
			var d dto.CommunityDiscussionResponse
			if err := discRows.Scan(&d.ID, &d.Author, &d.Avatar, &d.Content, &d.CreatedAt); err == nil {
				discussions = append(discussions, d)
			}
		}
		discRows.Close()
		detail.Discussions = discussions
	}

	return &detail, nil
}

func (r *CommunityRepository) Join(ctx context.Context, communityID int, userID int) error {
	query := `
		INSERT INTO community_members (community_id, user_id, community_role, joined_at)
		VALUES ($1, $2, 'member', NOW())
		ON CONFLICT DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, communityID, userID)
	return err

}

func (r *CommunityRepository) Leave(ctx context.Context, communityID int, userID int) error {
	query := `DELETE FROM community_members WHERE community_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, communityID, userID)
	return err
}

func (r *CommunityRepository) AddDiscussion(ctx context.Context, communityID int, userID int, content string) error {
	query := `
		INSERT INTO community_discussions (community_id, user_id, content, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(ctx, query, communityID, userID, content)
	return err
}

func (r *CommunityRepository) GetPopular(ctx context.Context, limit int) ([]dto.PopularCommunityResponse, error) {
	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT 
			c.id,
			c.name,
			c.slug,
			COALESCE(c.description, ''),
			COALESCE(c.cover_image_url, ''),
			COUNT(DISTINCT cm.id) AS total_members,
			COUNT(DISTINCT e.id) FILTER (WHERE e.start_time > CURRENT_TIMESTAMP AND e.status = 'published') AS total_upcoming_events,
			(COUNT(DISTINCT cm.id) * 1) + 
			(COUNT(DISTINCT e.id) FILTER (WHERE e.start_time > CURRENT_TIMESTAMP AND e.status = 'published') * 5) + 
			(COALESCE(SUM(CASE WHEN er.status = 'registered' THEN 1 ELSE 0 END), 0) * 2) AS popularity_score
		FROM communities c
		LEFT JOIN community_members cm ON c.id = cm.community_id
		LEFT JOIN events e ON c.id = e.community_id
		LEFT JOIN event_registrations er ON e.id = er.event_id
		WHERE c.is_active = TRUE
		GROUP BY c.id
		ORDER BY popularity_score DESC, total_members DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.PopularCommunityResponse
	for rows.Next() {
		var p dto.PopularCommunityResponse
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.CoverImageURL,
			&p.TotalMembers,
			&p.TotalUpcomingEvents,
			&p.PopularityScore,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *CommunityRepository) GetMembers(ctx context.Context, communityID int) ([]dto.CommunityMemberDetail, error) {
	query := `
		SELECT 
			u.id,
			u.full_name,
			COALESCE(u.avatar_url, ''),
			COALESCE(u.location, ''),
			COALESCE(cm.community_role, 'member'),
			cm.joined_at
		FROM community_members cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.community_id = $1
		ORDER BY cm.joined_at ASC
	`

	rows, err := r.db.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []dto.CommunityMemberDetail
	for rows.Next() {
		var m dto.CommunityMemberDetail
		if err := rows.Scan(
			&m.ID,
			&m.FullName,
			&m.AvatarURL,
			&m.Location,
			&m.CommunityRole,
			&m.JoinedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}
