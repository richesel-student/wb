package model

import "time"

type Event struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Date           time.Time     `json:"date"`
	Capacity       int           `json:"capacity"`
	AvailableSeats int           `json:"available_seats"`
	BookingTTL     time.Duration `json:"booking_ttl_minutes"`
}
