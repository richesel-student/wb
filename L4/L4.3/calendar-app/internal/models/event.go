package models

import "time"

// Event — модель события
type Event struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ReminderBeforeMinutes int       `json:"reminder_before_minutes"`
	Archived             bool      `json:"archived"`
}