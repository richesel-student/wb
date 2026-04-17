package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func New(client *minio.Client, bucket string) *Storage {
	return &Storage{client: client, bucket: bucket}
}

func (s *Storage) Upload(ctx context.Context, key string, r io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{})
	return err
}

func (s *Storage) Download(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return io.ReadAll(obj)
}

func (s *Storage) UploadBytes(ctx context.Context, key string, data []byte) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}
