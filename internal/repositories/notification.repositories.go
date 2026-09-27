package repositories

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetMyNotifications(ctx context.Context, userID int) ([]dto.NotificationResponse, error) {
	query := `
		SELECT 
			id,
			title,
			description,
			type,
			COALESCE(reference_type, ''),
			COALESCE(reference_id, 0),
			is_read,
			created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.NotificationResponse
	for rows.Next() {
		var n dto.NotificationResponse
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Description,
			&n.Type,
			&n.ReferenceType,
			&n.ReferenceID,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	return result, nil
}
