package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"shortener/internal/usecase"
)

type Handler struct {
	uc *usecase.ShortenerUseCase
}

func NewHandler(uc *usecase.ShortenerUseCase) *Handler {
	return &Handler{uc}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL    string `json:"url"`
		Custom string `json:"custom"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	short, err := h.uc.Create(req.URL, req.Custom)
	if err != nil {
		http.Error(w, "failed to create short url", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"short": short,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
} // 👈 ВАЖНО

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	short := strings.TrimPrefix(r.URL.Path, "/s/")

	original, err := h.uc.Redirect(short, r.UserAgent())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func (h *Handler) Analytics(w http.ResponseWriter, r *http.Request) {
	short := strings.TrimPrefix(r.URL.Path, "/analytics/")

	data, err := h.uc.Analytics(short)
	if err != nil {
		http.Error(w, "failed to get analytics", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
