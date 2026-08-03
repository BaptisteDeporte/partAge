package files

import (
	"baptistedeporte/partage/internal/storage"
	"context"
	"crypto/rand"
	"io"
)

type Storage interface {
	Put(ctx context.Context, id string, r io.Reader) error
	Get(ctx context.Context, id string) (io.ReadSeekCloser, *storage.FileMetadata, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Upload(ctx context.Context, r io.Reader) (string, error) {
	id := rand.Text()
	if err := s.storage.Put(ctx, id, r); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) Download(ctx context.Context, id string) (io.ReadSeekCloser, *storage.FileMetadata, error) {
	rs, fi, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return rs, fi, nil
}
