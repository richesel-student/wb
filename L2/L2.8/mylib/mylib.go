package mylib

import (
	"github.com/beevik/ntp"
	"time"
)

// TimeNTP получает точное время с NTP-сервера.
func TimeNTP() (time.Time, error) {
	return ntp.Time("pool.ntp.org")
}

