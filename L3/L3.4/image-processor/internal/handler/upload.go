package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"image-processor/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, id, key string, r io.Reader) error
}

type Repo interface {
	Get(ctx context.Context, id string) (*models.Image, error)
	GetAll(ctx context.Context) ([]*models.Image, error)
	Delete(ctx context.Context, id string) error
}

type Storage interface {
	Delete(ctx context.Context, key string) error
}

type Handler struct {
	svc  Service
	repo Repo
	st   Storage
}

func New(svc Service, repo Repo, st Storage) *Handler {
	return &Handler{svc: svc, repo: repo, st: st}
}

type ImageResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Original  string `json:"original"`
	Processed string `json:"processed,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func mapToResponse(img *models.Image) *ImageResponse {
	resp := &ImageResponse{
		ID:       img.ID,
		Status:   string(img.Status),
		Original: img.OriginalPath,
	}

	if img.ProcessedPath != "" {
		resp.Processed = img.ProcessedPath
		resp.Thumbnail = "thumbnails/" + img.ID
	}

	return resp
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "file is required", map[string]string{
			"file": "missing multipart field",
		})
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		writeJSONError(w, http.StatusBadRequest, "failed to read file", nil)
		return
	}

	contentType := http.DetectContentType(buf[:n])
	if contentType != "image/jpeg" && contentType != "image/png" {
		writeJSONError(w, http.StatusBadRequest, "invalid file type", map[string]string{
			"file": "only JPEG and PNG allowed",
		})
		return
	}

	id := uuid.NewString()
	key := "original/" + id

	reader := io.MultiReader(bytes.NewReader(buf[:n]), file)
	if err := h.svc.Upload(r.Context(), id, key, reader); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GET /image/{id}
func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/image/")

	img, err := h.repo.Get(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found", nil)
		return
	}

	json.NewEncoder(w).Encode(mapToResponse(img))
}

// DELETE /image/{id}
func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/image/")

	img, err := h.repo.Get(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found", nil)
		return
	}

	_ = h.st.Delete(r.Context(), img.OriginalPath)
	if img.ProcessedPath != "" {
		_ = h.st.Delete(r.Context(), img.ProcessedPath)
		_ = h.st.Delete(r.Context(), "thumbnails/"+img.ID)
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /images
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var resp []*ImageResponse
	for _, img := range list {
		resp = append(resp, mapToResponse(img))
	}

	json.NewEncoder(w).Encode(resp)
}
