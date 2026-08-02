package files

import (
	"context"
	"crypto/rand"
	"io"
)

type Storage interface {
	Put(ctx context.Context, id string, r io.Reader) error
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
