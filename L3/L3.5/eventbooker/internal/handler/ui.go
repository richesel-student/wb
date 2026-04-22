package handler

import (
	"embed"
	"eventbooker/internal/service"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templateFS embed.FS

type Handler struct {
	svc *service.Service
	tpl *template.Template
}

func New(s *service.Service) *Handler {
	tpl := template.Must(template.ParseFS(templateFS, "templates/index.html"))
	return &Handler{svc: s, tpl: tpl}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	events, _ := h.svc.ListEvents(r.Context())
	h.tpl.Execute(w, events)
}
