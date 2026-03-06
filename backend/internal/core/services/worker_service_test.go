package services

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkerRepository struct {
	mock.Mock
}

func (m *MockWorkerRepository) Get(ctx context.Context, userID, id string) (*domain.Worker, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Worker), args.Error(1)
}

func (m *MockWorkerRepository) List(ctx context.Context, userID, siteID string) ([]domain.Worker, error) {
	args := m.Called(ctx, userID, siteID)
	return args.Get(0).([]domain.Worker), args.Error(1)
}

func (m *MockWorkerRepository) Create(ctx context.Context, w *domain.Worker) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWorkerRepository) Update(ctx context.Context, w *domain.Worker) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWorkerRepository) Delete(ctx context.Context, userID, id string) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

func (m *MockWorkerRepository) ListByIsSynced(ctx context.Context, userID string, syncStatus int) ([]domain.Worker, error) {
	args := m.Called(ctx, userID, syncStatus)
	return args.Get(0).([]domain.Worker), args.Error(1)
}

func (m *MockWorkerRepository) GetProjectUserID(ctx context.Context, projectID string) (string, error) {
	args := m.Called(ctx, projectID)
	return args.String(0), args.Error(1)
}

func (m *MockWorkerRepository) GetByFIN(ctx context.Context, fin string) (*domain.Worker, error) {
	args := m.Called(ctx, fin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Worker), args.Error(1)
}

func (m *MockWorkerRepository) MarkSynced(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWorkerRepository) AssignToProject(ctx context.Context, projectID string, workerIDs []string, userID string) error {
	args := m.Called(ctx, projectID, workerIDs, userID)
	return args.Error(0)
}

func TestWorkerService_CreateWorker_Validation(t *testing.T) {
	mockRepo := new(MockWorkerRepository)
	svc := NewWorkerService(mockRepo)
	ctx := context.Background()

	// Invalid NRIC
	w := &domain.Worker{
		Name:       "Test Worker",
		PersonIDNo: "INVALID",
	}
	err := svc.CreateWorker(ctx, w)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid person_id_no")

	// Valid worker
	w2 := &domain.Worker{
		UserID:     "user1",
		Name:       "John Doe",
		PersonIDNo: "S1234567A",
		Status:     "active",
	}
	mockRepo.On("Create", ctx, w2).Return(nil)
	err = svc.CreateWorker(ctx, w2)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWorkerService_UpdateWorker_SyncTrigger(t *testing.T) {
	mockRepo := new(MockWorkerRepository)
	svc := NewWorkerService(mockRepo)
	ctx := context.Background()

	existing := &domain.Worker{
		ID:         "w1",
		UserID:     "user1",
		Name:       "John",
		FaceImgLoc: "old.jpg",
		IsSynced:   domain.SyncStatusSynced,
	}

	mockRepo.On("Get", ctx, "user1", "w1").Return(existing, nil)

	// Updating name should trigger sync
	newName := "John Updated"
	req := &domain.UpdateWorkerRequest{
		Name: &newName,
	}

	updated := *existing
	updated.Name = newName
	updated.IsSynced = domain.SyncStatusPendingUpdate

	mockRepo.On("Update", ctx, &updated).Return(nil)

	err := svc.UpdateWorker(ctx, "user1", "w1", req)
	assert.NoError(t, err)
	assert.Equal(t, domain.SyncStatusPendingUpdate, updated.IsSynced)
	mockRepo.AssertExpectations(t)
}
