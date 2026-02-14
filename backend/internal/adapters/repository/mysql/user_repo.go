package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT user_id, user_name, username, user_type, status, latitude, longitude FROM users WHERE username = ? AND status = 'active'`

	var t domain.User
	var lat, lng sql.NullString
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&t.ID, &t.Name, &t.Username, &t.UserType, &t.Status, &lat, &lng,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	if lat.Valid {
		t.Latitude = lat.String
	}
	if lng.Valid {
		t.Longitude = lng.String
	}

	return &t, nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT user_id, user_name, username, user_type, status, latitude, longitude FROM users WHERE user_id = ?`

	var t domain.User
	var lat, lng sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Username, &t.UserType, &t.Status, &lat, &lng,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lat.Valid {
		t.Latitude = lat.String
	}
	if lng.Valid {
		t.Longitude = lng.String
	}

	return &t, nil
}
