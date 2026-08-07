package files_test

import (
	"baptistedeporte/partage/internal/storage"
	"bytes"
	"context"
	"io"
	"time"
)

type fakeStorage struct {
	puts map[string][]byte
	err  error
}

func (f *fakeStorage) Put(ctx context.Context, id string, r io.Reader) error {
	if f.err != nil {
		return f.err
	}
	data, _ := io.ReadAll(r)
	f.puts[id] = data
	return nil
}

func (f *fakeStorage) Get(ctx context.Context, id string) (io.ReadSeekCloser, storage.FileMetadata, error) {
	data, ok := f.puts[id]
	if !ok {
		return nil, storage.FileMetadata{}, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)).(io.ReadSeekCloser),
		storage.FileMetadata{Name: id, LastMod: time.Now()},
		nil
}
