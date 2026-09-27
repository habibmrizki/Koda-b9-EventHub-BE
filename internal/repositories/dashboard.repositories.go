package repositories

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetAdminDashboard(ctx context.Context) (dto.AdminDashboardResponse, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM users) AS total_users,
			(SELECT COUNT(*) FROM events) AS total_events,
			(SELECT COUNT(*) FROM events WHERE status = 'published' AND start_time > CURRENT_TIMESTAMP) AS active_events,
			(SELECT COUNT(*) FROM communities WHERE is_active = TRUE) AS total_communities,
			(
				SELECT COALESCE(ROUND((COUNT(er.id) FILTER (WHERE er.status = 'registered')::NUMERIC / NULLIF(SUM(e.capacity), 0)::NUMERIC) * 100, 1), 0)
				FROM events e
				LEFT JOIN event_registrations er ON e.id = er.event_id
			) AS global_avg_fill_rate;
	`

	var res dto.AdminDashboardResponse
	err := r.db.QueryRow(ctx, query).Scan(
		&res.TotalUsers,
		&res.TotalEvents,
		&res.ActiveEvents,
		&res.TotalCommunities,
		&res.GlobalAvgFillRate,
	)
	if err != nil {
		return dto.AdminDashboardResponse{}, err
	}

	return res, nil
}
