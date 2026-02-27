package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
	"sgbuildex/internal/pkg/timeutil"
	"time"
)

type WorkerService struct {
	repo ports.WorkerRepository
}

func NewWorkerService(repo ports.WorkerRepository) ports.WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) GetWorker(ctx context.Context, userID, id string) (*domain.Worker, error) {
	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id scope required")
	}
	worker, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

func (s *WorkerService) ListWorkers(ctx context.Context, userID, siteID string) ([]domain.Worker, error) {
	return s.repo.List(ctx, userID, siteID)
}

func (s *WorkerService) CreateWorker(ctx context.Context, w *domain.Worker) error {
	if w.ID == "" {
		w.ID = "w" + time.Now().Format("20060102150405")
	}

	if w.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Default user_type if not set
	if w.UserType == "" {
		w.UserType = "user"
	}

	if w.CurrentProjectID != "" {
		projectUserID, err := s.repo.GetProjectUserID(ctx, w.CurrentProjectID)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}
		if projectUserID != w.UserID {
			return fmt.Errorf("user mismatch: project belongs to %s, worker belongs to %s", projectUserID, w.UserID)
		}
	}

	// If worker has auth data, mark as pending registration
	hasAuthData := w.FaceImgLoc != "" || w.CardNumber != ""
	if hasAuthData {
		w.IsSynced = domain.SyncStatusPendingRegistration
	}

	return s.repo.Create(ctx, w)
}

func (s *WorkerService) UpdateWorker(ctx context.Context, userID, id string, payload map[string]interface{}) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}

	existing, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return err
	}

	requiresSync := false
	wasRegistered := existing.FaceImgLoc != "" || existing.CardNumber != ""

	// Dynamic overlay logic
	if v, ok := payload["name"].(string); ok {
		if existing.Name != v {
			requiresSync = true
		}
		existing.Name = v
	}
	if v, ok := payload["email"].(string); ok {
		existing.Email = v
	}
	if v, ok := payload["role"].(string); ok {
		existing.Role = v
	}
	if v, ok := payload["user_type"].(string); ok {
		existing.UserType = v
	}
	if v, ok := payload["status"].(string); ok {
		existing.Status = v
		if v == "inactive" {
			existing.CurrentProjectID = ""
		}
	}
	if v, ok := payload["person_id_no"].(string); ok {
		existing.PersonIDNo = v
	}
	if v, ok := payload["user_id"].(string); ok {
		existing.UserID = v
	}
	if v, ok := payload["person_id_and_work_pass_type"].(string); ok {
		existing.PersonIDAndWorkPassType = v
	}
	if v, ok := payload["person_nationality"].(string); ok {
		existing.PersonNationality = v
	}
	if v, ok := payload["person_trade"].(string); ok {
		existing.PersonTrade = v
	}

	if v, ok := payload["auth_start_time"].(string); ok {
		if timeutil.CleanDateTime(existing.AuthStartTime) != timeutil.CleanDateTime(v) {
			requiresSync = true
		}
		existing.AuthStartTime = v
	}
	if v, ok := payload["auth_end_time"].(string); ok {
		if timeutil.CleanDateTime(existing.AuthEndTime) != timeutil.CleanDateTime(v) {
			requiresSync = true
		}
		existing.AuthEndTime = v
	}
	if v, ok := payload["face_img_loc"].(string); ok {
		if existing.FaceImgLoc != v {
			requiresSync = true
		}
		existing.FaceImgLoc = v
	}
	if v, ok := payload["card_number"].(string); ok {
		if existing.CardNumber != v {
			requiresSync = true
		}
		existing.CardNumber = v
	}
	if v, ok := payload["card_type"].(string); ok {
		if existing.CardType != v {
			requiresSync = true
		}
		existing.CardType = v
	}
	if v, ok := payload["fdid"].(float64); ok {
		if existing.FDID != int(v) {
			requiresSync = true
		}
		existing.FDID = int(v)
	}
	if v, ok := payload["fdid"].(int); ok {
		if existing.FDID != v {
			requiresSync = true
		}
		existing.FDID = v
	}

	isRegistered := existing.FaceImgLoc != "" || existing.CardNumber != ""

	// Project ID flexibility
	var newProjectID string
	if v, ok := payload["current_project_id"].(string); ok {
		newProjectID = v
	} else if v, ok := payload["project_id"].(string); ok {
		newProjectID = v
	} else if v, ok := payload["current_project_id"].(float64); ok {
		newProjectID = fmt.Sprintf("%.0f", v)
	} else {
		newProjectID = existing.CurrentProjectID
	}

	if newProjectID != existing.CurrentProjectID {
		if newProjectID != "" {
			projectUserID, err := s.repo.GetProjectUserID(ctx, newProjectID)
			if err != nil {
				return fmt.Errorf("invalid project ID: %w", err)
			}
			if projectUserID != existing.UserID {
				return fmt.Errorf("security violation: cannot assign project from different user")
			}
		}
		existing.CurrentProjectID = newProjectID
		requiresSync = true
	} else if v, ok := payload["current_project_id"]; ok && v == nil {
		if existing.CurrentProjectID != "" {
			requiresSync = true
		}
		existing.CurrentProjectID = ""
	}

	if requiresSync && (wasRegistered || isRegistered) {
		// Only set to pending update if it was previously synced or already pending update
		if existing.IsSynced != domain.SyncStatusPendingRegistration {
			existing.IsSynced = domain.SyncStatusPendingUpdate
		}
	} else {
		if v, ok := payload["is_synced"].(float64); ok {
			existing.IsSynced = int(v)
		} else if v, ok := payload["is_synced"].(int); ok {
			existing.IsSynced = v
		}
	}

	return s.repo.Update(ctx, existing)
}

func (s *WorkerService) DeleteWorker(ctx context.Context, userID, id string) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Delete(ctx, userID, id)
}

func (s *WorkerService) ListPendingSyncWorkers(ctx context.Context, userID string) ([]domain.Worker, error) {
	// Fetch workers needing registration and update
	registerWorkers, err := s.repo.ListByIsSynced(ctx, userID, domain.SyncStatusPendingRegistration)
	if err != nil {
		return nil, fmt.Errorf("failed to list register-pending workers: %w", err)
	}

	updateWorkers, err := s.repo.ListByIsSynced(ctx, userID, domain.SyncStatusPendingUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to list update-pending workers: %w", err)
	}

	return append(registerWorkers, updateWorkers...), nil
}

func (s *WorkerService) AssignWorkersToProject(ctx context.Context, projectID string, workerIDs []string) error {
	userID, err := s.repo.GetProjectUserID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to verify project: %w", err)
	}

	return s.repo.AssignToProject(ctx, projectID, workerIDs, userID)
}
