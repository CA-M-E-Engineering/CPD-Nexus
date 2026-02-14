package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type WorkerService struct {
	repo ports.WorkerRepository
}

func NewWorkerService(repo ports.WorkerRepository) ports.WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) GetWorker(ctx context.Context, id string) (*domain.Worker, error) {
	return s.repo.Get(ctx, id)
}

func (s *WorkerService) ListWorkers(ctx context.Context, tenantID, siteID string) ([]domain.Worker, error) {
	return s.repo.List(ctx, tenantID, siteID)
}

func (s *WorkerService) CreateWorker(ctx context.Context, w *domain.Worker) error {
	if w.ID == "" {
		w.ID = "user-" + uuid.New().String()
	}

	if w.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	if w.CurrentProjectID != "" {
		projectTenantID, err := s.repo.GetProjectTenantID(ctx, w.CurrentProjectID)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}
		if projectTenantID != w.TenantID {
			return fmt.Errorf("tenant mismatch: project belongs to %s, worker belongs to %s", projectTenantID, w.TenantID)
		}
	}

	return s.repo.Create(ctx, w)
}

func (s *WorkerService) UpdateWorker(ctx context.Context, id string, payload map[string]interface{}) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("worker not found")
	}

	// Dynamic overlay logic
	if v, ok := payload["name"].(string); ok {
		existing.Name = v
	}
	if v, ok := payload["email"].(string); ok {
		existing.Email = v
	}
	if v, ok := payload["role"].(string); ok {
		existing.Role = v
	}
	if v, ok := payload["status"].(string); ok {
		existing.Status = v
	}
	if v, ok := payload["trade_code"].(string); ok {
		existing.TradeCode = v
	}
	if v, ok := payload["fin"].(string); ok {
		existing.FIN = v
	}
	if v, ok := payload["company_name"].(string); ok {
		existing.CompanyName = v
	}
	if v, ok := payload["tenant_id"].(string); ok {
		existing.TenantID = v
	}

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

	if newProjectID != "" && newProjectID != existing.CurrentProjectID {
		projectTenantID, err := s.repo.GetProjectTenantID(ctx, newProjectID)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}
		if projectTenantID != existing.TenantID {
			return fmt.Errorf("security violation: cannot assign project from different tenant")
		}
		existing.CurrentProjectID = newProjectID
	} else if v, ok := payload["current_project_id"]; ok && v == nil {
		existing.CurrentProjectID = ""
	}

	return s.repo.Update(ctx, existing)
}

func (s *WorkerService) DeleteWorker(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
