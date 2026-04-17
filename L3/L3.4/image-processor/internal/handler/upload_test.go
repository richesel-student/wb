package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"image-processor/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct{ mock.Mock }

func (m *MockService) Upload(ctx context.Context, id, key string, r io.Reader) error {
	args := m.Called(ctx, id, key, r)
	return args.Error(0)
}

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Get(ctx context.Context, id string) (*models.Image, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Image), args.Error(1)
}

func (m *MockRepo) GetAll(ctx context.Context) ([]*models.Image, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Image), args.Error(1)
}

func (m *MockRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockRepo) UpdateStatus(ctx context.Context, id, status, path string) error {
	return m.Called(ctx, id, status, path).Error(0)
}

type MockStorage struct{ mock.Mock }

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}

func createMultipartJPEG() (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, _ := w.CreateFormFile("file", "img.jpg")

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	part.Write(buf.Bytes())

	w.Close()
	return body, w.FormDataContentType()
}

func TestUpload_Success(t *testing.T) {
	svc := new(MockService)
	repo := new(MockRepo)
	st := new(MockStorage)

	svc.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	h := New(svc, repo, st)

	body, ct := createMultipartJPEG()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", ct)

	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpload_InvalidType(t *testing.T) {
	h := New(new(MockService), new(MockRepo), new(MockStorage))

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", "file.txt")
	part.Write([]byte("not image"))
	w.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetImage(t *testing.T) {
	repo := new(MockRepo)

	repo.On("Get", mock.Anything, "1").Return(&models.Image{
		ID:            "1",
		Status:        models.StatusDone,
		OriginalPath:  "original/1",
		ProcessedPath: "processed/1",
	}, nil)

	h := New(nil, repo, nil)

	req := httptest.NewRequest("GET", "/image/1", nil)
	rec := httptest.NewRecorder()

	h.GetImage(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.Equal(t, "1", resp["id"])
}
