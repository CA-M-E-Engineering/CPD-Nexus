package handlers

import (
	"context"
	"fmt"
	"log"
	"sgbuildex/internal/bridge"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"strconv"
	"time"
)

// UserSyncPayload matches the outbound REGISTER_USER / UPDATE_USER structure
type UserSyncPayload struct {
	Devices []string     `json:"devices"`
	User    UserSyncData `json:"user"`
}

type UserSyncData struct {
	EmployeeNo     string       `json:"employee_no"`
	Name           string       `json:"name"`
	UserType       string       `json:"user_type"`
	Validity       UserValidity `json:"validity"`
	Authentication UserAuth     `json:"authentication"`
}

type UserValidity struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UserAuth struct {
	Card *UserCard `json:"card,omitempty"`
	Face *UserFace `json:"face,omitempty"`
}

type UserCard struct {
	CardNo   string `json:"card_no"`
	CardType string `json:"card_type"`
}

type UserFace struct {
	FaceID  string `json:"face_id"`
	FaceURL string `json:"face_url"`
}

// UserSyncBuilder builds outbound REGISTER_USER and UPDATE_USER bridge messages
type UserSyncBuilder struct {
	workerService ports.WorkerService
	workerRepo    ports.WorkerRepository
	deviceRepo    ports.DeviceRepository
}

func NewUserSyncBuilder(
	workerService ports.WorkerService,
	workerRepo ports.WorkerRepository,
	deviceRepo ports.DeviceRepository,
) *UserSyncBuilder {
	return &UserSyncBuilder{
		workerService: workerService,
		workerRepo:    workerRepo,
		deviceRepo:    deviceRepo,
	}
}

// BuildSyncRequests finds all pending-sync workers and builds bridge messages for them.
// If userID is provided, only workers belonging to that user are processed.
// Returns:
// - messages: payloads to send
// - processedIDs: IDs of workers that will be marked as status 1
// - invalidWorkers: workers with no site/devices
// - unauthWorkers: workers with no face/card data
func (b *UserSyncBuilder) BuildSyncRequests(ctx context.Context, userID string) ([]bridge.Message, []string, []domain.Worker, []domain.Worker, error) {
	workers, err := b.workerService.ListPendingSyncWorkers(ctx, userID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to list pending sync workers: %w", err)
	}

	if len(workers) == 0 {
		return nil, nil, nil, nil, nil
	}

	log.Printf("[UserSync] Found %d workers pending sync", len(workers))

	var messages []bridge.Message
	var processedWorkerIDs []string
	var invalidWorkers []domain.Worker
	var unauthWorkers []domain.Worker

	for _, w := range workers {
		// 1. Check for biometric/card data first
		hasAuth := w.FaceImgLoc != "" || w.CardNumber != ""
		if !hasAuth {
			log.Printf("[UserSync] Worker %s (%s) has no face/card data, skipping", w.ID, w.Name)
			unauthWorkers = append(unauthWorkers, w)
			continue
		}

		// 2. Determine which site's devices to target
		if w.SiteID == "" {
			log.Printf("[UserSync] Worker %s (%s) has no site via project, skipping", w.ID, w.Name)
			invalidWorkers = append(invalidWorkers, w)
			continue
		}

		// 3. Get device SNs for the worker's site
		deviceSNs, err := b.deviceRepo.ListSNsBySiteID(ctx, w.SiteID)
		if err != nil {
			log.Printf("[UserSync] Failed to get devices for site %s: %v", w.SiteID, err)
			invalidWorkers = append(invalidWorkers, w)
			continue
		}

		if len(deviceSNs) == 0 {
			log.Printf("[UserSync] No devices found for site %s (worker %s), skipping", w.SiteID, w.ID)
			invalidWorkers = append(invalidWorkers, w)
			continue
		}

		// Determine action based on is_synced value
		action := "UPDATE_USER"
		if w.IsSynced == 2 {
			action = "REGISTER_USER"
		}

		// Format auth times
		startTime := formatSyncTime(w.AuthStartTime)
		endTime := formatSyncTime(w.AuthEndTime)

		// Build payload
		payload := UserSyncPayload{
			Devices: deviceSNs,
			User: UserSyncData{
				EmployeeNo: w.PersonIDNo,
				Name:       w.Name,
				UserType:   w.UserType,
				Validity: UserValidity{
					StartTime: startTime,
					EndTime:   endTime,
				},
			},
		}

		// Add card authentication if present
		if w.CardNumber != "" {
			cardType := w.CardType
			if cardType == "" {
				cardType = "normal"
			}
			payload.User.Authentication.Card = &UserCard{
				CardNo:   w.CardNumber,
				CardType: cardType,
			}
		}

		// Add face authentication if present
		if w.FaceImgLoc != "" {
			payload.User.Authentication.Face = &UserFace{
				FaceID:  strconv.Itoa(w.FDID),
				FaceURL: w.FaceImgLoc,
			}
		}

		// Build the bridge message
		msg, err := bridge.NewRequest(action, payload)
		if err != nil {
			log.Printf("[UserSync] Failed to build %s request for worker %s: %v", action, w.ID, err)
			continue
		}

		log.Printf("[UserSync] Built %s request for worker %s (%s) → %d devices at site %s",
			action, w.ID, w.Name, len(deviceSNs), w.SiteID)

		messages = append(messages, msg)
		processedWorkerIDs = append(processedWorkerIDs, w.ID)
	}

	return messages, processedWorkerIDs, invalidWorkers, unauthWorkers, nil
}

// MarkWorkersSynced marks the given workers as synced (is_synced=1) after successful send
func (b *UserSyncBuilder) MarkWorkersSynced(ctx context.Context, workerIDs []string) {
	for _, id := range workerIDs {
		err := b.workerService.UpdateWorker(ctx, id, map[string]interface{}{
			"is_synced": 1,
		})
		if err != nil {
			log.Printf("[UserSync] Failed to mark worker %s as synced: %v", id, err)
		}
	}
}

// formatSyncTime converts DB datetime to RFC3339 format for bridge payload
func formatSyncTime(t string) string {
	if t == "" {
		return ""
	}

	// Try parsing common formats
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, t); err == nil {
			return parsed.Format(time.RFC3339)
		}
	}

	// Return as-is if no format matches
	return t
}
