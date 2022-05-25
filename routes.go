package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/sokudo/filesystems/miniofilesystem"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)

	a.get("/test-minio", func(w http.ResponseWriter, r *http.Request) {
		f := a.App.FileSystems["MINIO"].(miniofilesystem.Minio)

		files, err := f.List("")
		if err != nil {
			log.Println(err)
			return
		}

		for _, file := range files {
			log.Println(file.Key)
		}
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
