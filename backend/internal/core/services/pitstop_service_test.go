package services

import (
	"context"
	"errors"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockPitstopRepository struct {
	mock.Mock
}

func (m *MockPitstopRepository) InsertAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error {
	args := m.Called(ctx, auths)
	return args.Error(0)
}

func (m *MockPitstopRepository) UpdateAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error {
	args := m.Called(ctx, auths)
	return args.Error(0)
}

func (m *MockPitstopRepository) GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.PitstopAuthorisation), args.Error(1)
}

func (m *MockPitstopRepository) AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error {
	args := m.Called(ctx, userID, onBehalfOfNames)
	return args.Error(0)
}

type MockAttendanceRepository struct {
	mock.Mock
}

func (m *MockAttendanceRepository) Get(ctx context.Context, userID, id string) (*domain.Attendance, error) {
	args := m.Called(ctx, userID, id)
	return args.Get(0).(*domain.Attendance), args.Error(1)
}

func (m *MockAttendanceRepository) List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	args := m.Called(ctx, userID, siteID, workerID, date)
	return args.Get(0).([]domain.Attendance), args.Error(1)
}

func (m *MockAttendanceRepository) ExtractProjectsWithPendingAttendance(ctx context.Context, userID string) ([]domain.Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockAttendanceRepository) ExtractPendingAttendance(ctx context.Context) ([]domain.AttendanceRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.AttendanceRow), args.Error(1)
}

func (m *MockAttendanceRepository) ExtractPendingAttendanceByProject(ctx context.Context, userID, projectID string) ([]domain.AttendanceRow, error) {
	args := m.Called(ctx, userID, projectID)
	return args.Get(0).([]domain.AttendanceRow), args.Error(1)
}

func (m *MockAttendanceRepository) Create(ctx context.Context, a *domain.Attendance) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAttendanceRepository) GetMaxID(ctx context.Context, pattern string) (string, error) {
	args := m.Called(ctx, pattern)
	return args.String(0), args.Error(1)
}

type MockSubmissionRepository struct {
	mock.Mock
}

func (m *MockSubmissionRepository) UpdateAttendanceStatus(ctx context.Context, attendanceID, status, payload, message string) error {
	args := m.Called(ctx, attendanceID, status, payload, message)
	return args.Error(0)
}

func (m *MockSubmissionRepository) LogSubmission(ctx context.Context, refType, refID, status, payload, message string) error {
	args := m.Called(ctx, refType, refID, status, payload, message)
	return args.Error(0)
}

type MockSettingsRepository struct {
	mock.Mock
}

func (m *MockSettingsRepository) GetSettings(ctx context.Context) (*domain.SystemSettings, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.SystemSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSettingsRepository) UpdateSettings(ctx context.Context, settings domain.SystemSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockSettingsRepository) GetDeviceStats(ctx context.Context) (total int, online int, err error) {
	args := m.Called(ctx)
	return args.Int(0), args.Int(1), args.Error(2)
}

type MockExternalSubmitter struct {
	mock.Mock
}

func (m *MockExternalSubmitter) SubmitManpowerUtilization(ctx context.Context, repo ports.SubmissionRepository, settings *domain.SystemSettings, rows []domain.AttendanceRow) (int, int, error) {
	args := m.Called(ctx, repo, settings, rows)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockExternalSubmitter) FetchPitstopConfig(ctx context.Context) (*ports.PitstopConfigResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*ports.PitstopConfigResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestPitstopService_TestSubmission_SuccessAndFailures(t *testing.T) {
	mockPitstopRepo := new(MockPitstopRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockSettingsRepo := new(MockSettingsRepository)
	mockExternalSubmitter := new(MockExternalSubmitter)

	svc := NewPitstopService(mockPitstopRepo, mockExternalSubmitter, mockAttendanceRepo, mockSubmissionRepo, mockSettingsRepo)

	ctx := context.Background()
	userID := "user123"
	projectID := "proj123"

	// Mock settings
	settings := &domain.SystemSettings{
		MaxWorkersPerRequest: 100,
		MaxPayloadSizeKB:     256,
		MaxRequestsPerMinute: 150,
	}
	mockSettingsRepo.On("GetSettings", ctx).Return(settings, nil)

	// Mock attendance rows (1 valid, 1 invalid missing FIN)
	rows := []domain.AttendanceRow{
		{
			AttendanceID:       "ATT-1",
			WorkerFIN:          "G1234567P",
			WorkerWorkPassType: "WP",
			WorkerTrade:        "2.3",
			EmployerName:       "Test Employer",
			EmployerUEN:        "12345678A",
			RegulatorID:        "REG-1",
			OnBehalfOfID:       "OB-1",
			SubmissionEntity:   1,
		},
		{
			AttendanceID: "ATT-2",
			WorkerFIN:    "", // missing FIN makes this row invalid
		},
	}
	mockAttendanceRepo.On("ExtractPendingAttendanceByProject", ctx, userID, projectID).Return(rows, nil)

	// Since ATT-2 is invalid, the adapter will update it as failed in the DB.
	// We mock the adapter to return 1 submitted and 1 failed.

	// The submitter will be called with the rows (internally filtered by the mapper)
	// Let's assume the external submitter successfully processes the 1 valid payload.
	mockExternalSubmitter.On("SubmitManpowerUtilization", ctx, mockSubmissionRepo, settings, rows).Return(1, 1, nil)

	// Execute TestSubmission
	submittedCount, failedCount, err := svc.TestSubmission(ctx, userID, projectID)

	// Verify expectations
	assert.NoError(t, err)
	assert.Equal(t, 1, submittedCount)
	assert.Equal(t, 1, failedCount)

	mockSettingsRepo.AssertExpectations(t)
	mockAttendanceRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
	mockExternalSubmitter.AssertExpectations(t)
}

func TestPitstopService_TestSubmission_ExtractError(t *testing.T) {
	mockPitstopRepo := new(MockPitstopRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockSettingsRepo := new(MockSettingsRepository)
	mockExternalSubmitter := new(MockExternalSubmitter)

	svc := NewPitstopService(mockPitstopRepo, mockExternalSubmitter, mockAttendanceRepo, mockSubmissionRepo, mockSettingsRepo)

	ctx := context.Background()

	mockSettingsRepo.On("GetSettings", ctx).Return(&domain.SystemSettings{}, nil)

	expectedErr := errors.New("db error")
	mockAttendanceRepo.On("ExtractPendingAttendanceByProject", ctx, "user123", "proj123").Return([]domain.AttendanceRow{}, expectedErr)

	submittedCount, failedCount, err := svc.TestSubmission(ctx, "user123", "proj123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract project attendance")
	assert.Equal(t, 0, submittedCount)
	assert.Equal(t, 0, failedCount)
}

func TestPitstopService_TestSubmission_NoRows(t *testing.T) {
	mockPitstopRepo := new(MockPitstopRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockSettingsRepo := new(MockSettingsRepository)
	mockExternalSubmitter := new(MockExternalSubmitter)

	svc := NewPitstopService(mockPitstopRepo, mockExternalSubmitter, mockAttendanceRepo, mockSubmissionRepo, mockSettingsRepo)

	ctx := context.Background()

	mockSettingsRepo.On("GetSettings", ctx).Return(&domain.SystemSettings{}, nil)
	mockAttendanceRepo.On("ExtractPendingAttendanceByProject", ctx, "user123", "proj123").Return([]domain.AttendanceRow{}, nil)

	submittedCount, failedCount, err := svc.TestSubmission(ctx, "user123", "proj123")

	assert.NoError(t, err)
	assert.Equal(t, 0, submittedCount)
	assert.Equal(t, 0, failedCount)
}
