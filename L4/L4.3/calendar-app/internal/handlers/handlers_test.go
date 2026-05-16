package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/models"
	"calendar-app/internal/repository"
)

// Вспомогательный метод для быстрой инициализации окружения
func setupTestEnv() (*Handler, *repository.Repository) {
	log := logger.NewAsyncLogger(10)
	log.Start(context.Background())
	repo := repository.NewRepository()
	h := NewHandler(repo, log)
	return h, repo
}

func TestHealthCheckHandler(t *testing.T) {
	h, _ := setupTestEnv()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(h.HealthCheck).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected 200, got %v", status)
	}

	expected := `{"status":"ok"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

func TestCreateEventHandler_Success(t *testing.T) {
	h, _ := setupTestEnv()

	reqBody, _ := json.Marshal(CreateEventRequest{
		Title:                 "Test Action",
		StartTime:             time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		EndTime:               time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		ReminderBeforeMinutes: 10,
	})

	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.CreateEvent).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 21, got %v", rr.Code)
	}
}

func TestCreateEventHandler_InvalidJSON(t *testing.T) {
	h, _ := setupTestEnv()

	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString("{invalid json"))
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.CreateEvent).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %v", rr.Code)
	}
}

func TestGetEventsHandler(t *testing.T) {
	h, repo := setupTestEnv()

	// Добавим одно событие в репозиторий напрямую
	_ = repo.Create(context.Background(), &models.Event{ID: "1", Title: "Event 1"})

	req, _ := http.NewRequest("GET", "/events", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.GetEvents).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", rr.Code)
	}
}

func TestGetEventHandler_NotFound(t *testing.T) {
	h, _ := setupTestEnv()

	// Запрашиваем несуществующий ID
	req, _ := http.NewRequest("GET", "/events/999", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.GetEvent).ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for missing event, got %v", rr.Code)
	}
}

func TestUpdateEventHandler_NotFound(t *testing.T) {
	h, _ := setupTestEnv()

	reqBody, _ := json.Marshal(CreateEventRequest{Title: "Updated"})
	req, _ := http.NewRequest("PUT", "/events/999", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.UpdateEvent).ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 on update, got %v", rr.Code)
	}
}

func TestDeleteEventHandler_NotFound(t *testing.T) {
	h, _ := setupTestEnv()

	req, _ := http.NewRequest("DELETE", "/events/999", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(h.DeleteEvent).ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 on delete, got %v", rr.Code)
	}
}