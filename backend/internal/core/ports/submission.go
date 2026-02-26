package ports

import "context"

type SubmissionRepository interface {
	LogSubmission(ctx context.Context, dataElementID, internalID, status, payload, errorMessage string) error
	UpdateAttendanceStatus(ctx context.Context, attendanceID, status, responsePayload, errorMessage string) error
}
