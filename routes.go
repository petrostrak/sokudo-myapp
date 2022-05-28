package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)

	a.get("/users/login", a.Handlers.UserLogin)
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.Logout)

	a.get("/auth/{provider}", a.Handlers.SocialLogin)
	a.get("/auth/{provider}/callback", a.Handlers.SocialMediaCallback)

	a.get("/upload", a.Handlers.SokudoUpload)
	a.post("/upload", a.Handlers.PostSokudoUpload)

	a.get("/list-fs", a.Handlers.ListFS)

	a.get("/files/upload", a.Handlers.UploadToFS)
	a.post("/files/upload", a.Handlers.PostUploadToFS)

	a.get("/delete-from-fs", a.Handlers.DeleteFromFS)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
