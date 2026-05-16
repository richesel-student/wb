package worker

import (
	"context"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/models"
)

type ReminderTask struct {
	EventID  string
	Title    string
	RemindAt time.Time
}

type ReminderWorker struct {
	tasks  chan ReminderTask
	logger *logger.AsyncLogger
}

func NewReminderWorker(l *logger.AsyncLogger) *ReminderWorker {
	return &ReminderWorker{
		tasks:  make(chan ReminderTask, 100),
		logger: l,
	}
}

func (w *ReminderWorker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.logger.Info("Reminder worker stopped")
				return
			case task := <-w.tasks:
				w.processTask(ctx, task)
			}
		}
	}()
}

func (w *ReminderWorker) processTask(ctx context.Context, task ReminderTask) {
	now := time.Now()
	if task.RemindAt.After(now) {
		sleepDuration := task.RemindAt.Sub(now)
		timer := time.NewTimer(sleepDuration)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			w.sendReminder(task)
		}
	} else {
		// Если время напоминания уже прошло, отправляем немедленно
		w.sendReminder(task)
	}
}

func (w *ReminderWorker) sendReminder(task ReminderTask) {
	w.logger.Info("REMINDER: Event is scheduled for now or soon!", map[string]interface{}{"title": task.Title, "id": task.EventID})
}

func (w *ReminderWorker) ScheduleReminder(event *models.Event) {
	if event.ReminderBeforeMinutes <= 0 {
		return
	}

	remindAt := event.StartTime.Add(-time.Duration(event.ReminderBeforeMinutes) * time.Minute)

	task := ReminderTask{
		EventID:  event.ID,
		Title:    event.Title,
		RemindAt: remindAt,
	}

	select {
	case w.tasks <- task:
		w.logger.Info("Reminder scheduled for event", map[string]interface{}{"title": event.Title, "remindAt": remindAt})
	default:
		w.logger.Error("Reminder channel is full, skipping reminder", map[string]interface{}{"title": event.Title})
	}
}

// GetTasksChannel возвращает канал для тестирования
func (w *ReminderWorker) GetTasksChannel() chan ReminderTask {
	return w.tasks
}