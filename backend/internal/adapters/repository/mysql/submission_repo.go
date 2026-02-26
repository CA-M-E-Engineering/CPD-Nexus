package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/ports"
	"time"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) ports.SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) LogSubmission(ctx context.Context, dataElementID, internalID, status, payload, errorMessage string) error {
	query := `
		INSERT INTO submission_logs (data_element_id, internal_id, status, payload, error_message)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, dataElementID, internalID, status, payload, errorMessage)
	return err
}

func (r *SubmissionRepository) UpdateAttendanceStatus(ctx context.Context, attendanceID, status, responsePayload, errorMessage string) error {
	query := `
		UPDATE attendance
		SET status = ?, response_payload = ?, error_message = ?, updated_at = ?
		WHERE attendance_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, responsePayload, errorMessage, time.Now(), attendanceID)
	return err
}
