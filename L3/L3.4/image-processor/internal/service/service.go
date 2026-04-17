package service

import (
	"context"
	"io"

	"image-processor/internal/models"
)

type Repo interface {
	Create(ctx context.Context, img *models.Image) error
}

type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64) error
}

type Queue interface {
	Send(ctx context.Context, id string) error
}

type Service struct {
	repo Repo
	st   Storage
	q    Queue
}

func New(repo Repo, storage Storage, queue Queue) *Service {
	return &Service{repo: repo, st: storage, q: queue}
}

func (s *Service) Upload(ctx context.Context, id, key string, r io.Reader) error {
	if err := s.st.Upload(ctx, key, r, -1); err != nil {
		return err
	}
	img := &models.Image{
		ID:           id,
		Status:       models.StatusQueued,
		OriginalPath: key,
	}
	if err := s.repo.Create(ctx, img); err != nil {
		return err
	}
	return s.q.Send(ctx, id)
}
