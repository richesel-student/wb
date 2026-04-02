package service

import (
	"time"

	"delayed-notifier/internal/model"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
)

type NotificationService struct{}

func (s *NotificationService) Create(n model.Notification) (string, error) {
	n.Status = "pending"
	n.CreatedAt = time.Now()

	if err := repository.Create(n); err != nil {
		return "", err
	}

	delay := time.Until(n.SendAt)
	if n.SendAt.IsZero() {
		delay = 5 * time.Second
	}

	if err := queue.Publish(n, delay); err != nil {
		return "", err
	}

	return n.ID, nil
}

func (s *NotificationService) Cancel(id string) error {
	if err := repository.UpdateStatus(id, "canceled"); err != nil {
		return err
	}
	return nil
}
