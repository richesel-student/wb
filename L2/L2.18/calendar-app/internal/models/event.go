package models

import "time"

// Event — модель события
type Event struct {
	UserID int       // пользователь
	Date   time.Time // дата события
	Text   string    // описание
}