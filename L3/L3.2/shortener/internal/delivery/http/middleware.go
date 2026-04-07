package http

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// =====================
// RATE LIMIT
// =====================

var limiter = rate.NewLimiter(10, 20) // 10 req/sec, burst 20

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// =====================
// LOGGING
// =====================

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.UserAgent(),
			time.Since(start),
		)
	}
}

// =====================
// CHAIN
// =====================

// позволяет комбинировать middleware
func Chain(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
