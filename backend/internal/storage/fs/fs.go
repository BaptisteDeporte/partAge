package fs

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrAlreadyExists = errors.New("file already exists")
	ErrNotADirectory = errors.New("not a directory")
	ErrInvalidID     = errors.New("invalid file id")
)

type FS struct {
	rootDir string
}

func New(dir string) (*FS, error) {
	i, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !i.IsDir() {
		return nil, ErrNotADirectory
	}

	fs := &FS{
		rootDir: dir,
	}

	return fs, nil
}

func (f *FS) Put(ctx context.Context, id string, r io.Reader) (err error) {
	if strings.ContainsAny(id, `/\`) || id == "" || id == "." || id == ".." {
		return ErrInvalidID
	}

	p := filepath.Join(f.rootDir, id)
	rel, err := filepath.Rel(f.rootDir, p)

	if err != nil || strings.HasPrefix(rel, "..") {
		return ErrInvalidID
	}

	if _, err := os.Stat(p); err == nil {
		return ErrAlreadyExists
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	tmp, err := os.CreateTemp(f.rootDir, "upload-*")
	if err != nil {
		return err
	}

	committed := false

	defer func() {
		if !committed {
			os.Remove(tmp.Name())
		}
	}()

	defer func() {
		if cerr := tmp.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(tmp, r); err != nil {
		return err
	}

	if err = tmp.Sync(); err != nil {
		return err
	}

	if err = os.Rename(tmp.Name(), p); err != nil {
		return err
	}

	committed = true

	return
}
