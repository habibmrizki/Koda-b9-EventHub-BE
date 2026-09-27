package repositories

import (
	"context"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestimonialRepository struct {
	db *pgxpool.Pool
}

func NewTestimonialRepository(db *pgxpool.Pool) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

func (r *TestimonialRepository) GetFeatured(ctx context.Context, limit int) ([]dto.TestimonialResponse, error) {
	if limit <= 0 {
		limit = 6
	}

	query := `
		SELECT 
			t.id,
			t.content,
			COALESCE(t.author_title, ''),
			t.created_at,
			u.full_name,
			COALESCE(u.avatar_url, ''),
			t.is_featured
		FROM testimonials t
		JOIN users u ON t.user_id = u.id
		WHERE t.is_featured = TRUE
		ORDER BY t.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []dto.TestimonialResponse
	for rows.Next() {
		var item dto.TestimonialResponse
		if err := rows.Scan(
			&item.ID,
			&item.Content,
			&item.AuthorTitle,
			&item.CreatedAt,
			&item.FullName,
			&item.AvatarURL,
			&item.IsFeatured,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *TestimonialRepository) GetAll(ctx context.Context) ([]dto.TestimonialResponse, error) {
	query := `
		SELECT 
			t.id,
			t.content,
			COALESCE(t.author_title, ''),
			t.created_at,
			u.full_name,
			COALESCE(u.avatar_url, ''),
			t.is_featured
		FROM testimonials t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.TestimonialResponse
	for rows.Next() {
		var item dto.TestimonialResponse
		if err := rows.Scan(
			&item.ID,
			&item.Content,
			&item.AuthorTitle,
			&item.CreatedAt,
			&item.FullName,
			&item.AvatarURL,
			&item.IsFeatured,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *TestimonialRepository) Create(ctx context.Context, userID int, req dto.CreateTestimonialRequest) (dto.TestimonialResponse, error) {
	query := `
		WITH inserted AS (
			INSERT INTO testimonials (user_id, content, author_title, is_featured, created_at)
			VALUES ($1, $2, $3, TRUE, CURRENT_TIMESTAMP)
			RETURNING id, content, author_title, is_featured, created_at
		)
		SELECT 
			i.id,
			i.content,
			i.author_title,
			u.full_name,
			COALESCE(u.avatar_url, ''),
			i.is_featured,
			i.created_at
		FROM inserted i
		JOIN users u ON u.id = $1
	`

	var item dto.TestimonialResponse
	err := r.db.QueryRow(ctx, query, userID, req.Content, req.AuthorTitle).Scan(
		&item.ID,
		&item.Content,
		&item.AuthorTitle,
		&item.FullName,
		&item.AvatarURL,
		&item.IsFeatured,
		&item.CreatedAt,
	)
	if err != nil {
		return dto.TestimonialResponse{}, err
	}

	return item, nil
}
