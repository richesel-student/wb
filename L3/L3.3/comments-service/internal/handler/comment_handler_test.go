package handler

import (
	"bytes"
	"comments-service/internal/repository"
	"comments-service/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupHandler() *CommentHandler {
	repo := repository.NewMemoryRepository()
	svc := service.NewCommentService(repo)
	return NewCommentHandler(svc)
}

func TestCreateHandler(t *testing.T) {
	h := setupHandler()

	body := []byte(`{"text":"hello"}`)

	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetHandler(t *testing.T) {
	h := setupHandler()

	// сначала создаём
	createReq := httptest.NewRequest(http.MethodPost, "/comments",
		bytes.NewBuffer([]byte(`{"text":"hello"}`)))
	createW := httptest.NewRecorder()
	h.Create(createW, createReq)

	req := httptest.NewRequest(http.MethodGet, "/comments", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Body.String() == "null" {
		t.Fatal("expected non-null response")
	}
}

func TestDeleteHandler(t *testing.T) {
	h := setupHandler()

	// create
	createReq := httptest.NewRequest(http.MethodPost, "/comments",
		bytes.NewBuffer([]byte(`{"text":"hello"}`)))
	createW := httptest.NewRecorder()
	h.Create(createW, createReq)

	req := httptest.NewRequest(http.MethodDelete, "/comments/1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
