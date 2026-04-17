package service

import (
	"bytes"
	"context"
	"io"
	"testing"

	"image-processor/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Create(ctx context.Context, img *models.Image) error {
	return m.Called(ctx, img).Error(0)
}

type MockStorage struct{ mock.Mock }

func (m *MockStorage) Upload(ctx context.Context, key string, r io.Reader, size int64) error {
	return m.Called(ctx, key, r, size).Error(0)
}

type MockQueue struct{ mock.Mock }

func (m *MockQueue) Send(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestUpload_Flow(t *testing.T) {
	repo := new(MockRepo)
	st := new(MockStorage)
	q := new(MockQueue)

	svc := New(repo, st, q)

	st.On("Upload", mock.Anything, "key", mock.Anything, int64(-1)).Return(nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Image")).Return(nil)
	q.On("Send", mock.Anything, "id").Return(nil)

	err := svc.Upload(context.Background(), "id", "key", bytes.NewReader([]byte("data")))

	assert.NoError(t, err)

	st.AssertCalled(t, "Upload", mock.Anything, "key", mock.Anything, int64(-1))
	repo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*models.Image"))
	q.AssertCalled(t, "Send", mock.Anything, "id")
}

func TestUpload_StorageFail(t *testing.T) {
	st := new(MockStorage)
	st.On("Upload", mock.Anything, mock.Anything, mock.Anything, int64(-1)).Return(assert.AnError)

	svc := New(new(MockRepo), st, new(MockQueue))

	err := svc.Upload(context.Background(), "id", "key", bytes.NewReader([]byte("data")))

	assert.Error(t, err)
}
