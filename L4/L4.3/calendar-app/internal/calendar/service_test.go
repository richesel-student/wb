package calendar

import (
	"testing"
	"time"
)

// TestCreateEvent проверяет создание события
func TestCreateEvent(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "test"); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForDay(1, date)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

// TestUpdateEvent проверяет обновление существующего события
func TestUpdateEvent(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "old"); err != nil {
		t.Fatal(err)
	}

	if err := s.UpdateEvent(1, date, "new"); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForDay(1, date)

	if events[0].Title != "new" {
		t.Fatalf("expected 'new', got '%s'", events[0].Title)
	}
}

// TestUpdateEvent_NotFound проверяет ошибку при обновлении несуществующего события
func TestUpdateEvent_NotFound(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	err := s.UpdateEvent(1, date, "new")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestDeleteEvent проверяет удаление события
func TestDeleteEvent(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "test"); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteEvent(1, date); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForDay(1, date)

	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

// TestDeleteEvent_NotFound проверяет ошибку при удалении несуществующего события
func TestDeleteEvent_NotFound(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	err := s.DeleteEvent(1, date)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestEventsForDay проверяет получение событий за день
func TestEventsForDay(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "event1"); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForDay(1, date)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

// TestEventsForDay_Empty проверяет случай без событий
func TestEventsForDay_Empty(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	events, _ := s.EventsForDay(1, date)

	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

// TestEventsForWeek проверяет получение событий за неделю
func TestEventsForWeek(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "event1"); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForWeek(1, date)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

// TestEventsForMonth проверяет получение событий за месяц
func TestEventsForMonth(t *testing.T) {
	s := NewService()

	date, _ := time.Parse("2006-01-02", "2025-03-28")

	if err := s.CreateEvent(1, date, "event1"); err != nil {
		t.Fatal(err)
	}

	events, _ := s.EventsForMonth(1, date)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}