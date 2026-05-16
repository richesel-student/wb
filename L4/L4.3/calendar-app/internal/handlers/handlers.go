package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"calendar-app/internal/logger"
	"calendar-app/internal/models"
	"calendar-app/internal/repository"
	"calendar-app/internal/worker"
)

// Handler — HTTP слой
type Handler struct {
	repo           *repository.Repository
	reminderWorker *worker.ReminderWorker
	logger         *logger.AsyncLogger
}

func NewHandler(repo *repository.Repository, l *logger.AsyncLogger) *Handler {
	return &Handler{
		repo:           repo,
		reminderWorker: worker.NewReminderWorker(l),
		logger:         l,
	}
}

// SetReminderWorker устанавливает reminder worker
func (h *Handler) SetReminderWorker(worker *worker.ReminderWorker) {
	h.reminderWorker = worker
}

// ===== ВСПОМОГАТЕЛЬНОЕ =====

// CreateEventRequest — запрос на создание события
type CreateEventRequest struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	StartTime             string `json:"start_time"`
	EndTime               string `json:"end_time"`
	ReminderBeforeMinutes int    `json:"reminder_before_minutes"`
}

// UpdateEventRequest — запрос на обновление события
type UpdateEventRequest struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	StartTime             string `json:"start_time"`
	EndTime               string `json:"end_time"`
	ReminderBeforeMinutes int    `json:"reminder_before_minutes"`
}

// writeJSON — универсальный ответ
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// parseTime — парсинг времени ISO 8601
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// ===== HANDLERS =====

// CreateEvent — POST /events
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	startTime, err := parseTime(req.StartTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid start_time format"})
		return
	}

	endTime, err := parseTime(req.EndTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid end_time format"})
		return
	}

	if endTime.Before(startTime) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "end_time must be after start_time"})
		return
	}

	event := &models.Event{
		ID:                    fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:                 req.Title,
		Description:           req.Description,
		StartTime:             startTime,
		EndTime:               endTime,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		ReminderBeforeMinutes: req.ReminderBeforeMinutes,
		Archived:              false,
	}

	if err := h.repo.Create(r.Context(), event); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Планирование напоминания
	h.reminderWorker.ScheduleReminder(event)

	writeJSON(w, http.StatusCreated, event)
}

// GetEvents — GET /events
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	events, err := h.repo.GetActive(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, events)
}

// GetEvent — GET /events/{id}
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/events/")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "event id is required"})
		return
	}

	event, ok, err := h.repo.Get(r.Context(), path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	writeJSON(w, http.StatusOK, event)
}

// UpdateEvent — PUT /events/{id}
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/events/")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "event id is required"})
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	existingEvent, ok, err := h.repo.Get(r.Context(), path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	startTime, err := parseTime(req.StartTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid start_time format"})
		return
	}

	endTime, err := parseTime(req.EndTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid end_time format"})
		return
	}

	if endTime.Before(startTime) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "end_time must be after start_time"})
		return
	}

	updatedEvent := &models.Event{
		ID:                    path,
		Title:                 req.Title,
		Description:           req.Description,
		StartTime:             startTime,
		EndTime:               endTime,
		CreatedAt:             existingEvent.CreatedAt,
		UpdatedAt:             time.Now(),
		ReminderBeforeMinutes: req.ReminderBeforeMinutes,
		Archived:              existingEvent.Archived,
	}

	ok, err = h.repo.Update(r.Context(), path, updatedEvent)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "failed to update event: not found"})
		return
	}

	// Обновление напоминания, если изменилось время или настройки напоминания
	if existingEvent.ReminderBeforeMinutes > 0 || updatedEvent.ReminderBeforeMinutes > 0 {
		h.reminderWorker.ScheduleReminder(updatedEvent)
	}

	writeJSON(w, http.StatusOK, updatedEvent)
}

// DeleteEvent — DELETE /events/{id}
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/events/")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "event id is required"})
		return
	}

	ok, err := h.repo.Delete(r.Context(), path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "event deleted successfully"})
}

// HealthCheck — GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}