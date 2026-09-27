package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, userID int) (models.Users, error) {
	query := `
		SELECT id, full_name, email, password, role, location, COALESCE(avatar_url, ''), bio, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var user models.Users
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Location,
		&user.AvatarURL,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.Users{}, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID int, req dto.UpdateProfileRequest) (models.Users, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.FullName != nil {
		setClauses = append(setClauses, fmt.Sprintf("full_name = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.FullName))
		argPos++
	}

	if req.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("location = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Location))
		argPos++
	}

	if req.AvatarURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_url = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.AvatarURL))
		argPos++
	}

	if req.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argPos))
		args = append(args, strings.TrimSpace(*req.Bio))
		argPos++
	}

	if len(setClauses) == 0 {
		return r.FindByID(ctx, userID)
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, full_name, email, password, role, location, COALESCE(avatar_url, ''), bio, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPos)
	args = append(args, userID)

	var user models.Users
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Location,
		&user.AvatarURL,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	res, err := r.db.Exec(ctx, query, hashedPassword, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
