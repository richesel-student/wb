package worker

import (
	"context"
	"testing"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/models"
)

func TestReminderWorker_ScheduleReminder(t *testing.T) {
	log := logger.NewAsyncLogger(10)
	log.Start(context.Background())

	w := NewReminderWorker(log)

	event := &models.Event{
		ID:                    "test-id",
		Title:                 "Meeting",
		StartTime:             time.Now().Add(30 * time.Minute),
		ReminderBeforeMinutes: 10,
	}

	w.ScheduleReminder(event)

	select {
	case task := <-w.GetTasksChannel():
		if task.EventID != event.ID {
			t.Errorf("expected %s, got %s", event.ID, task.EventID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for reminder task")
	}
}

func TestReminderWorker_ImmediateReminder(t *testing.T) {
	log := logger.NewAsyncLogger(10)
	log.Start(context.Background())

	w := NewReminderWorker(log)

	// Событие уже идет, напоминание должно сработать мгновенно
	event := &models.Event{
		ID:                    "immediate-id",
		Title:                 "Urgent",
		StartTime:             time.Now().Add(-5 * time.Minute),
		ReminderBeforeMinutes: 10,
	}

	w.ScheduleReminder(event)

	select {
	case task := <-w.GetTasksChannel():
		if task.EventID != event.ID {
			t.Errorf("expected %s, got %s", event.ID, task.EventID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for immediate task")
	}
}