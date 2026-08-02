package main

import (
	"log"
	"net/http"

	"baptistedeporte/partage/internal/files"
	phttp "baptistedeporte/partage/internal/http"
	"baptistedeporte/partage/internal/storage/fs"
)

func main() {
	st, err := fs.New("/tmp/partage")
	if err != nil {
		panic(err)
	}

	svc := files.New(st)

	router := phttp.Router(svc)

	log.Printf("listening on :3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("server: %v", err)
	}
}
