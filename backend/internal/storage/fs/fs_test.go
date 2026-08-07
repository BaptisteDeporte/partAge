package fs

import (
	"baptistedeporte/partage/internal/storage"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T) *FS {
	t.Helper()
	dir := t.TempDir()
	f, err := New(dir)
	if err != nil {
		t.Fatalf("failed to setup FS, %v", err)
	}
	return f
}

func TestGetPath(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"valid id", "abc123", nil},
		{"empty id", "", storage.ErrInvalidID},
		{"dot id", ".", storage.ErrInvalidID},
		{"double dot id", "..", storage.ErrInvalidID},
		{"contains slash", "foo/bar", storage.ErrInvalidID},
		{"path traversal", "../../etc/passwd", storage.ErrInvalidID},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := setup(t)
			_, err := f.getPath(tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestPutGetRoundtrip(t *testing.T) {
	f := setup(t)

	ctx := context.Background()
	const id = "testing-id"
	const content = "hello world"
	ri := strings.NewReader(content)

	if err := f.Put(ctx, id, ri); err != nil {
		t.Fatalf("failed to put file, %v", err)
	}

	ro, fm, err := f.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get file, %v", err)
	}

	defer ro.Close()

	if fm.Name != id {
		t.Errorf("metadata name: got %q, want %q", fm.Name, id)
	}
	if fm.LastMod.IsZero() {
		t.Errorf("metadata LastMod is zero")
	}

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, ro); err != nil {
		t.Fatalf("failed to check file content, %v", err)
	}

	if buf.String() != content {
		t.Fatalf("failed to get file, got %v, want %v", buf.String(), content)
	}
}

func TestPut_ErrAlreadyExists(t *testing.T) {
	f := setup(t)

	ctx := context.Background()
	const id = "testing-id"
	const content = "hello world"

	if err := f.Put(ctx, id, strings.NewReader(content)); err != nil {
		t.Fatalf("failed to put file, %v", err)
	}

	err := f.Put(ctx, id, strings.NewReader(content))

	if !errors.Is(err, storage.ErrAlreadyExists) {
		t.Fatalf("failed to trigger ErrAlreadyExists, %v", err)
	}
}

func TestGet_ErrNotFound(t *testing.T) {
	f := setup(t)

	ctx := context.Background()
	const id = "testing-id"

	_, _, err := f.Get(ctx, id)

	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("failed to trigger ErrNotFound, %v", err)
	}
}

type failingReader struct {
	data []byte
	pos  int
}

var errSimulated = errors.New("simulated failure")

func (r *failingReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data)/2 {
		return 0, errSimulated
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func TestPut_ErrFailedReader(t *testing.T) {
	f := setup(t)

	ctx := context.Background()
	const id = "testing-id"
	const content = "hello world"

	fr := failingReader{
		data: []byte(content),
	}

	if err := f.Put(ctx, id, &fr); !errors.Is(err, errSimulated) {
		t.Fatalf("failed to trigger ErrSimulated, %v", err)
	}

	if _, _, err := f.Get(ctx, id); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("failed to trigger ErrNotFound, %v", err)
	}

	entries, err := os.ReadDir(f.rootDir)
	if err != nil {
		t.Fatalf("read rootDir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected clean rootDir, got %d entries: %v", len(entries), entries)
	}
}
