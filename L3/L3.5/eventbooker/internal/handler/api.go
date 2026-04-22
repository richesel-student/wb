package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"eventbooker/internal/model"
)

type EventResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Capacity       int       `json:"capacity"`
	AvailableSeats int       `json:"available_seats"`
	BookingTTL     int       `json:"booking_ttl_minutes"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	capacity, err := strconv.Atoi(r.FormValue("capacity"))
	if err != nil || capacity <= 0 {
		http.Error(w, "bad capacity", 400)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "bad date", 400)
		return
	}

	ttlMin, _ := strconv.Atoi(r.FormValue("ttl"))
	ttl := time.Duration(ttlMin) * time.Minute

	id, err := h.svc.CreateEvent(r.Context(), model.Event{
		Name:       name,
		Date:       date,
		Capacity:   capacity,
		BookingTTL: ttl,
	})
	if err != nil {
		http.Error(w, "error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad id", 400)
		return
	}

	e, err := h.svc.GetEvent(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	resp := EventResponse{
		ID:             e.ID,
		Name:           e.Name,
		Date:           e.Date,
		Capacity:       e.Capacity,
		AvailableSeats: e.AvailableSeats,
		BookingTTL:     int(e.BookingTTL.Minutes()),
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Book(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad id", 400)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "missing user id", 401)
		return
	}

	event, err := h.svc.GetEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	bookingID, err := h.svc.Book(r.Context(), eventID, userID, event.BookingTTL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"booking_id": bookingID,
	})
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad id", 400)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "missing user id", 401)
		return
	}

	bookingID, err := h.svc.FindUserBooking(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, "booking not found", 404)
		return
	}

	if err := h.svc.Confirm(r.Context(), bookingID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
