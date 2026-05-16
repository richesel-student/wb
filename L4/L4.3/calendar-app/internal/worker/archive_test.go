package worker

import (
	"context"
	"testing"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/models"
	"calendar-app/internal/repository"
)

func TestArchiveWorker_ArchiveOldEvents(t *testing.T) {
	log := logger.NewAsyncLogger(10)
	log.Start(context.Background())

	repo := repository.NewRepository()
	w := NewArchiveWorker(repo, log)
	w.SetThreshold(0)

	expiredEvent := &models.Event{
		ID:        "expired-id",
		StartTime: time.Now().Add(-10 * time.Minute),
		EndTime:   time.Now().Add(-5 * time.Minute),
		Archived:  false,
	}

	_ = repo.Create(context.Background(), expiredEvent)
	w.archiveOldEvents()

	event, _, _ := repo.Get(context.Background(), "expired-id")
	if !event.Archived {
		t.Error("expected event to be archived")
	}
}

func TestArchiveWorker_DoNotArchiveFutureEvents(t *testing.T) {
	log := logger.NewAsyncLogger(10)
	log.Start(context.Background())

	repo := repository.NewRepository()
	w := NewArchiveWorker(repo, log)
	w.SetThreshold(24 * time.Hour)

	futureEvent := &models.Event{
		ID:        "future-id",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Archived:  false,
	}

	_ = repo.Create(context.Background(), futureEvent)
	w.archiveOldEvents()

	event, _, _ := repo.Get(context.Background(), "future-id")
	if event.Archived {
		t.Error("future event should NOT be archived")
	}
}