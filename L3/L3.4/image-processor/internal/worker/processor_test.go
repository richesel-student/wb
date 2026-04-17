package worker

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	"image-processor/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Get(ctx context.Context, id string) (*models.Image, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Image), args.Error(1)
}

func (m *MockRepo) UpdateStatus(ctx context.Context, id, status, path string) error {
	return m.Called(ctx, id, status, path).Error(0)
}

type MockStorage struct{ mock.Mock }

func (m *MockStorage) Download(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorage) UploadBytes(ctx context.Context, key string, data []byte) error {
	return m.Called(ctx, key, data).Error(0)
}

func createImageBytes() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.White)
		}
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

func TestHandle(t *testing.T) {
	repo := new(MockRepo)
	st := new(MockStorage)

	p := &Processor{repo: repo, st: st}

	imgData := createImageBytes()

	repo.On("Get", mock.Anything, "id").Return(&models.Image{
		ID:           "id",
		OriginalPath: "original/id",
	}, nil)

	repo.On("UpdateStatus", mock.Anything, "id", mock.Anything, mock.Anything).
		Return(nil).Maybe()

	st.On("Download", mock.Anything, "original/id").Return(imgData, nil)

	st.On("UploadBytes", mock.Anything, "processed/id.jpg", mock.AnythingOfType("[]uint8")).Return(nil)
	st.On("UploadBytes", mock.Anything, "thumbnails/id.jpg", mock.AnythingOfType("[]uint8")).Return(nil)

	err := p.handle(context.Background(), "id")

	assert.NoError(t, err)
}
