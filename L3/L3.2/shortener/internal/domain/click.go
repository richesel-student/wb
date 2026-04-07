package domain

import "time"

type Click struct {
	Short     string
	UserAgent string
	Time      time.Time
}
