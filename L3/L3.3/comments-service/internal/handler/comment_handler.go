package handler

import (
	"comments-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text     string `json:"text"`
		ParentID *int64 `json:"parent_id"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	comment, _ := h.service.Create(req.Text, req.ParentID)

	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) Get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var parentID *int64
	if p := query.Get("parent_id"); p != "" {
		id, _ := strconv.ParseInt(p, 10, 64)
		parentID = &id
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	search := query.Get("search")

	if limit == 0 {
		limit = 100
	}

	comments, _ := h.service.GetTree(parentID, limit, offset, search)
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comments/")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	h.service.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}
