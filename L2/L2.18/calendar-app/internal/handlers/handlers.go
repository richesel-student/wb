package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"calendar-app/internal/calendar"
)

// Handler — HTTP слой
type Handler struct {
	service *calendar.Service
}

func NewHandler(s *calendar.Service) *Handler {
	return &Handler{service: s}
}

// ===== ВСПОМОГАТЕЛЬНОЕ =====

// request — входящие данные
type request struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Event  string `json:"event"`
}

// parseRequest — поддержка JSON и form-data
func parseRequest(r *http.Request) (request, error) {
	var req request

	// JSON
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		return req, err
	}

	// form-data
	if err := r.ParseForm(); err != nil {
		return req, err
	}

	userID, _ := strconv.Atoi(r.FormValue("user_id"))

	req.UserID = userID
	req.Date = r.FormValue("date")
	req.Event = r.FormValue("event")

	return req, nil
}

// parseDate — парсинг даты YYYY-MM-DD
func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// writeJSON — универсальный ответ
func writeJSON(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ===== HANDLERS =====

// CREATE
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	date, err := parseDate(req.Date)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	err = h.service.CreateEvent(req.UserID, date, req.Event)
	if err != nil {
		writeJSON(w, 500, map[string]interface{}{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"result": "event created"})
}

// UPDATE
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	date, err := parseDate(req.Date)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	err = h.service.UpdateEvent(req.UserID, date, req.Event)
	if err != nil {
		writeJSON(w, 503, map[string]interface{}{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"result": "event updated"})
}

// DELETE
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	date, err := parseDate(req.Date)
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	err = h.service.DeleteEvent(req.UserID, date)
	if err != nil {
		writeJSON(w, 503, map[string]interface{}{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"result": "event deleted"})
}

// DAY
func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid user_id"})
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	events, _ := h.service.EventsForDay(userID, date)

	writeJSON(w, 200, map[string]interface{}{"result": events})
}

// WEEK
func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid user_id"})
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	events, _ := h.service.EventsForWeek(userID, date)

	writeJSON(w, 200, map[string]interface{}{"result": events})
}

// MONTH
func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid user_id"})
		return
	}

	date, err := parseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeJSON(w, 400, map[string]interface{}{"error": "invalid date"})
		return
	}

	events, _ := h.service.EventsForMonth(userID, date)

	writeJSON(w, 200, map[string]interface{}{"result": events})
}