package services

import (
	"context"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"time"
)

type ScheduledTask func(ctx context.Context)
type TimeExtractor func(settings *domain.SystemSettings) string

type DailyScheduler struct {
	settingsRepo  ports.SettingsRepository
	task          ScheduledTask
	timeExtractor TimeExtractor
	name          string
	reset         chan struct{}
}

func NewDailyScheduler(repo ports.SettingsRepository, name string, extractor TimeExtractor, task ScheduledTask) *DailyScheduler {
	return &DailyScheduler{
		settingsRepo:  repo,
		task:          task,
		timeExtractor: extractor,
		name:          name,
		reset:         make(chan struct{}, 1),
	}
}

func (s *DailyScheduler) Reset() {
	select {
	case s.reset <- struct{}{}:
	default: // Already reset pending
	}
}

// Start runs the scheduling loop.
func (s *DailyScheduler) Start(ctx context.Context) {
	for {
		// 1. Get the settings from database
		settings, err := s.settingsRepo.GetSettings(ctx)
		if err != nil {
			log.Printf("[%s] Scheduler: Failed to get settings: %v. Retrying in 1 minute...", s.name, err)
			select {
			case <-time.After(1 * time.Minute):
				continue
			case <-ctx.Done():
				return
			}
		}

		scheduledTimeStr := s.timeExtractor(settings) // Format should be "HH:MM:SS"
		now := time.Now()

		// Parse the HH:MM:SS
		var hour, min, sec int
		_, err = fmt.Sscanf(scheduledTimeStr, "%d:%d:%d", &hour, &min, &sec)
		if err != nil {
			log.Printf("[%s] Scheduler: Invalid time format '%s': %v. Retrying in 1 minute...", s.name, scheduledTimeStr, err)
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
		log.Printf("[%s] Scheduler: Next run scheduled for %v (in %v)", s.name, nextRun.Format(time.RFC3339), durationUntilNext.Truncate(time.Second))

		select {
		case <-time.After(durationUntilNext):
			log.Printf("[%s] Scheduler: [TRIGGER] Starting scheduled task...", s.name)
			s.task(ctx)
			log.Printf("[%s] Scheduler: [COMPLETED] task finished. Scheduling next run...", s.name)

			// Wait a few seconds to avoid double trigger
			select {
			case <-time.After(5 * time.Second):
			case <-s.reset: // Check for reset during the short wait too
			case <-ctx.Done():
				return
			}

		case <-s.reset:
			log.Printf("[%s] Scheduler: [RESET] Schedule updated, re-evaluating...", s.name)
			continue

		case <-ctx.Done():
			log.Printf("[%s] Scheduler: Shutting down...", s.name)
			return
		}
	}
}
