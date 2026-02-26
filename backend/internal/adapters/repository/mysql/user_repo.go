package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

const userBaseSelect = `
    SELECT 
        u.user_id, u.user_name, u.username, u.user_type, u.status, 
        u.latitude, u.longitude, u.contact_email, u.contact_phone, u.address, u.password_hash,
        (SELECT COUNT(*) FROM workers w WHERE w.user_id = u.user_id AND w.status = 'active') as worker_count,
        (SELECT COUNT(*) FROM devices d WHERE d.user_id = u.user_id AND d.status != 'inactive') as device_count
    FROM users u`

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := userBaseSelect + " WHERE u.username = ? AND u.status = ?"
	return r.scanRow(r.db.QueryRowContext(ctx, query, username, domain.StatusActive))
}

func (r *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	query := userBaseSelect + " WHERE u.user_id = ?"
	return r.scanRow(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, userBaseSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		u, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (
			user_id, user_name, user_type, contact_email, contact_phone, 
            username, password_hash, status, address, latitude, longitude
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Name, u.UserType, u.ContactEmail, u.ContactPhone,
		u.Username, u.PasswordHash, u.Status, u.Address, u.Latitude, u.Longitude)
	return err
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET 
			user_name=?, user_type=?, contact_email=?, contact_phone=?, username=?, 
            status=?, latitude=?, longitude=?, address=?, password_hash=?
		WHERE user_id=?`

	_, err := r.db.ExecContext(ctx, query,
		u.Name, u.UserType, u.ContactEmail, u.ContactPhone, u.Username,
		u.Status, u.Latitude, u.Longitude, u.Address, u.PasswordHash, u.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Deactivate user
	if _, err := tx.ExecContext(ctx, "UPDATE users SET status = ? WHERE user_id = ?", domain.StatusInactive, id); err != nil {
		return err
	}

	// 2. Deactivate all workers for this user
	if _, err := tx.ExecContext(ctx, "UPDATE workers SET status = ?, current_project_id = NULL WHERE user_id = ?", domain.StatusInactive, id); err != nil {
		return err
	}

	// 3. Deactivate all projects for this user
	if _, err := tx.ExecContext(ctx, "UPDATE projects SET status = 'inactive' WHERE user_id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) scanRow(scanner Scanner) (*domain.User, error) {
	var u domain.User
	var lat, lng sql.NullFloat64
	var email, phone, addr, hash sql.NullString

	err := scanner.Scan(
		&u.ID, &u.Name, &u.Username, &u.UserType, &u.Status,
		&lat, &lng, &email, &phone, &addr, &hash,
		&u.WorkerCount, &u.DeviceCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if lat.Valid {
		u.Latitude = lat.Float64
	}
	if lng.Valid {
		u.Longitude = lng.Float64
	}
	if email.Valid {
		u.ContactEmail = email.String
	}
	if phone.Valid {
		u.ContactPhone = phone.String
	}
	if addr.Valid {
		u.Address = addr.String
	}
	if hash.Valid {
		u.PasswordHash = hash.String
	}

	return &u, nil
}
