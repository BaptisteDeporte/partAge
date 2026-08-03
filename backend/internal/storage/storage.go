package storage

import (
	"errors"
	"time"
)

var (
	ErrAlreadyExists = errors.New("file already exists")
	ErrNotFound      = errors.New("not found")
	ErrInvalidID     = errors.New("invalid file id")
)

type FileMetadata struct {
	Name    string
	LastMod time.Time
}
