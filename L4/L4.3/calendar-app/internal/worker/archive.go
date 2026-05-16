package worker

import (
	"context"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/repository"
)

type ArchiveWorker struct {
	repo      *repository.Repository
	threshold time.Duration // события старше этого порога архивируются
	logger    *logger.AsyncLogger
}

func NewArchiveWorker(repo *repository.Repository, l *logger.AsyncLogger) *ArchiveWorker {
	return &ArchiveWorker{
		repo:      repo,
		threshold: 24 * time.Hour, // архивировать события старше 24 часов
		logger:    l,
	}
}

func (w *ArchiveWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	w.logger.Info("Archive worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Archive worker stopped")
			return
		case <-ticker.C:
			w.archiveOldEvents()
		}
	}
}

func (w *ArchiveWorker) archiveOldEvents() {
	events, err := w.repo.GetActive(context.Background())
	if err != nil {
		w.logger.Error("Failed to get active events for archiving", map[string]interface{}{"error": err.Error()})
		return
	}
	now := time.Now()
	archivedCount := 0

	for _, event := range events {
		// Архивируем события, которые закончились более threshold назад
		if event.EndTime.Add(w.threshold).Before(now) {
			if archived, err := w.repo.Archive(context.Background(), event.ID); err == nil && archived {
				archivedCount++
				w.logger.Info("Archived event", map[string]interface{}{"title": event.Title, "endTime": event.EndTime})
			} else if err != nil {
				w.logger.Error("Failed to archive event", map[string]interface{}{"error": err.Error()})
			}
		}
	}

	if archivedCount > 0 {
		w.logger.Info("Archived old events", map[string]interface{}{"count": archivedCount})
	}
}

// SetThreshold устанавливает порог для архивации
func (w *ArchiveWorker) SetThreshold(threshold time.Duration) {
	w.threshold = threshold
}