package calendar

import (
	"errors"
	"sync"
	"time"

	"calendar-app/internal/models"
)

// Service — бизнес-логика (НЕ зависит от HTTP)
type Service struct {
	mu     sync.RWMutex
	events map[int][]models.Event // userID → список событий
}

// NewService — конструктор
func NewService() *Service {
	return &Service{
		events: make(map[int][]models.Event),
	}
}

// CreateEvent — добавление нового события
func (s *Service) CreateEvent(userID int, date time.Time, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[userID] = append(s.events[userID], models.Event{
		UserID: userID,
		Date:   date,
		Text:   text,
	})

	return nil
}

// UpdateEvent — обновление события по дате
// (обновляет первое найденное событие за день)
func (s *Service) UpdateEvent(userID int, date time.Time, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.events[userID] {
		if sameDay(s.events[userID][i].Date, date) {
			s.events[userID][i].Text = text
			return nil
		}
	}

	return errors.New("event not found")
}

// DeleteEvent — удаляет ВСЕ события за указанный день
// (фикс бага: раньше удалялось только одно)
func (s *Service) DeleteEvent(userID int, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := s.events[userID]
	var filtered []models.Event

	found := false

	for _, e := range events {
		if sameDay(e.Date, date) {
			found = true
			continue
		}
		filtered = append(filtered, e)
	}

	if !found {
		return errors.New("event not found")
	}

	s.events[userID] = filtered
	return nil
}

// EventsForDay — получить события за день
func (s *Service) EventsForDay(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []models.Event

	for _, e := range s.events[userID] {
		if sameDay(e.Date, date) {
			res = append(res, e)
		}
	}

	return res, nil
}

// EventsForWeek — события за неделю (ISO неделя)
func (s *Service) EventsForWeek(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	year, week := date.ISOWeek()

	var res []models.Event

	for _, e := range s.events[userID] {
		y, w := e.Date.ISOWeek()
		if y == year && w == week {
			res = append(res, e)
		}
	}

	return res, nil
}

// EventsForMonth — события за месяц
func (s *Service) EventsForMonth(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []models.Event

	for _, e := range s.events[userID] {
		if e.Date.Year() == date.Year() &&
			e.Date.Month() == date.Month() {
			res = append(res, e)
		}
	}

	return res, nil
}

// sameDay — сравнение дат (игнорирует время)
func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}