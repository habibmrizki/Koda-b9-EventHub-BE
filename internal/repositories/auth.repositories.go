package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, fullName, email, hashedPassword string) (models.Users, error) {
	query := `
		INSERT INTO users (full_name, email, password, role, location, avatar_url, bio)
		VALUES ($1, $2, $3, 'attendee', '', '', '')
		RETURNING id, full_name, email, password, role, location, COALESCE(avatar_url, ''), bio, created_at, updated_at
	`

	var user models.Users
	err := r.db.QueryRow(ctx, query, fullName, email, hashedPassword).Scan(
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
		if strings.Contains(err.Error(), "users_email_key") || strings.Contains(err.Error(), "duplicate key") {
			return models.Users{}, ErrEmailExists
		}
		return models.Users{}, err
	}

	return user, nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (models.Users, error) {
	query := `
		SELECT id, full_name, email, password, role, location, COALESCE(avatar_url, ''), bio, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	var user models.Users
	err := r.db.QueryRow(ctx, query, email).Scan(
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
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Users{}, err
		}
		return models.Users{}, ErrUserNotFound
	}

	return user, nil
}

