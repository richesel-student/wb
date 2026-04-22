package handler

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// UI
	mux.HandleFunc("GET /{$}", h.Index)

	// API
	mux.HandleFunc("POST /api/events", h.CreateEvent)
	mux.HandleFunc("GET /api/events/{id}", h.GetEvent)
	mux.HandleFunc("POST /api/events/{id}/book", h.Book)
	mux.HandleFunc("POST /api/events/{id}/confirm", h.Confirm)

	mux.HandleFunc("GET /api/health", h.Health)

	return mux
}
