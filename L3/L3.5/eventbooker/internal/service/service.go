package service

import (
	"context"
	"time"

	"eventbooker/internal/model"
	"eventbooker/internal/repository"
)

type Service struct {
	repo repository.BookingRepository
}

func New(r repository.BookingRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateEvent(ctx context.Context, e model.Event) (int, error) {
	return s.repo.CreateEvent(ctx, e)
}

func (s *Service) GetEvent(ctx context.Context, id int) (model.Event, error) {
	return s.repo.GetEvent(ctx, id)
}

func (s *Service) ListEvents(ctx context.Context) ([]model.Event, error) {
	return s.repo.ListEvents(ctx)
}

func (s *Service) Book(ctx context.Context, eventID int, userID string, ttl time.Duration) (int, error) {
	return s.repo.CreateBooking(ctx, eventID, userID, ttl)
}

func (s *Service) Confirm(ctx context.Context, bookingID int) error {
	return s.repo.ConfirmBooking(ctx, bookingID)
}

func (s *Service) Cleanup(ctx context.Context) error {
	return s.repo.ExpireBookings(ctx, 1000)
}

func (s *Service) FindUserBooking(ctx context.Context, eventID int, userID string) (int, error) {
	return s.repo.FindUserBooking(ctx, eventID, userID)
}
