package services

import (
	"context"
	"fmt"
	"log"
	"sgbuildex/internal/core/ports"
	"time"
)

type SubmissionTask func(ctx context.Context)

type CPDScheduler struct {
	settingsRepo ports.SettingsRepository
	task         SubmissionTask
}

func NewCPDScheduler(repo ports.SettingsRepository, task SubmissionTask) *CPDScheduler {
	return &CPDScheduler{
		settingsRepo: repo,
		task:         task,
	}
}

// Start runs the scheduling loop. It fetches the submission time from the database
// and schedules the next task execution accordingly.
func (s *CPDScheduler) Start(ctx context.Context) {
	for {
		// 1. Get the submission time from settings
		settings, err := s.settingsRepo.GetSettings(ctx)
		if err != nil {
			log.Printf("Scheduler: Failed to get settings from database: %v. Retrying in 1 minute...", err)
			select {
			case <-time.After(1 * time.Minute):
				continue
			case <-ctx.Done():
				return
			}
		}

		submissionTimeStr := settings.CPDSubmissionTime // Format should be "HH:MM:SS"
		now := time.Now()

		// Parse the HH:MM:SS
		var hour, min, sec int
		_, err = fmt.Sscanf(submissionTimeStr, "%d:%d:%d", &hour, &min, &sec)
		if err != nil {
			log.Printf("Scheduler: Invalid submission time format '%s' in DB: %v. Retrying in 1 minute...", submissionTimeStr, err)
			select {
			case <-time.After(1 * time.Minute):
				continue
			case <-ctx.Done():
				return
			}
		}

		// Calculate when the next run should be (today)
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())

		// If it's already past this time today, schedule for tomorrow
		if !nextRun.After(now) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		durationUntilNext := time.Until(nextRun)
		log.Printf("Scheduler: Next CPD submission scheduled for %v (in %v)", nextRun.Format(time.RFC3339), durationUntilNext.Truncate(time.Second))

		select {
		case <-time.After(durationUntilNext):
			log.Println("Scheduler: [TRIGGER] Starting scheduled CPD submission...")
			s.task(ctx)
			log.Println("Scheduler: [COMPLETED] submission task finished. Scheduling next run...")

			// Wait a few seconds to ensure we don't accidentally re-trigger for the same second
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			log.Println("Scheduler: Shutting down...")
			return
		}
	}
}
