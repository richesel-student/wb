package model

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Channel   string    `json:"channel"`
	To        string    `json:"to"`
	SendAt    time.Time `json:"send_at"`
	Status    string    `json:"status"`
	Retries   int       `json:"retries"`
	CreatedAt time.Time `json:"created_at"`
}
