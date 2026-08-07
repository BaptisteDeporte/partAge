package fs

import (
	"baptistedeporte/partage/internal/storage"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotADirectory = errors.New("not a directory")

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

func (f *FS) getPath(id string) (string, error) {
	if strings.ContainsAny(id, `/\`) || id == "" || id == "." || id == ".." {
		return "", storage.ErrInvalidID
	}

	p := filepath.Join(f.rootDir, id)
	rel, err := filepath.Rel(f.rootDir, p)

	if err != nil || strings.HasPrefix(rel, "..") {
		return "", storage.ErrInvalidID
	}

	return p, nil
}

func (f *FS) Get(ctx context.Context, id string) (io.ReadSeekCloser, storage.FileMetadata, error) {
	p, err := f.getPath(id)
	if err != nil {
		return nil, storage.FileMetadata{}, err
	}

	file, err := os.Open(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, storage.FileMetadata{}, storage.ErrNotFound
		}
		return nil, storage.FileMetadata{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, storage.FileMetadata{}, err
	}

	fmeta := storage.FileMetadata{
		Name:    fileInfo.Name(),
		LastMod: fileInfo.ModTime(),
	}

	fmeta.Name = fileInfo.Name()
	fmeta.LastMod = fileInfo.ModTime()

	return file, fmeta, nil
}

func (f *FS) Put(ctx context.Context, id string, r io.Reader) (err error) {
	p, err := f.getPath(id)
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); err == nil {
		return storage.ErrAlreadyExists
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
