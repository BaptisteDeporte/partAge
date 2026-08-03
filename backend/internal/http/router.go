package http

import (
	"baptistedeporte/partage/internal/files"
	"net/http"
)

func Router(svc *files.Service) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("PUT /files/upload", PutFile(svc))
	router.HandleFunc("GET /files/{id}", GetFile(svc))

	return router
}
