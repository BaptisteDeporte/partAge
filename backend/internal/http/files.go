package http

import (
	"baptistedeporte/partage/internal/files"
	"baptistedeporte/partage/internal/storage"
	"errors"
	"log"
	"net/http"
)

const maxUploadSize = 100 << 20

func PutFile(svc *files.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		id, err := svc.Upload(ctx, r.Body)
		if err != nil {
			switch {
			case errors.As(err, new(*http.MaxBytesError)):
				http.Error(w, "too large", http.StatusRequestEntityTooLarge)
			default:
				log.Printf("upload failed: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(id))
	}
}

func GetFile(svc *files.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := r.Context()
		rs, fm, err := svc.Download(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				http.Error(w, "not found", http.StatusNotFound)
			default:
				log.Printf("download failed: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		defer rs.Close()

		http.ServeContent(w, r, fm.Name, fm.LastMod, rs)
	}
}
