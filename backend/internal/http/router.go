package http

import (
	"net/http"

	"baptistedeporte/partage/internal/files"
)

func Router(svc *files.Service) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("PUT /files/upload", PutFile(svc))

	return router
}
