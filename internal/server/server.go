package server

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates/*
var templatesFS embed.FS
var Templates = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

func NewEngine() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func NewServer(app *core.App) http.Handler {
	app.Templates = Templates
	loadViewerCSS()
	r := NewEngine()

	r.Use(middleware.Logger)
	r.Use(injectContentCSS)

	FileServer(r, "/static", "./static")
	FileServer(r, "/lib", "/opt/gochive/lib")
	FileServer(r, "/build", "/opt/gochive/lib/pdfjs/build/")
	FileServer(r, "/web", "/opt/gochive/lib/pdfjs/web/")

	r.Get("/", handlers.Root(app))
	r.Get("/file/{id}", handlers.GetFile(app))
	r.Get("/view/{id}", handlers.View(app))
	r.Get("/images/{id}", handlers.GetImage(app))
	r.Get("/get_files/{page}", handlers.GetFiles(app))

	r.Post("/search/{page}", handlers.SearchFiles(app))
	r.Post("/set_favorite/{id}", handlers.SetFavorite(app))
	r.Patch("/edit/{id}", handlers.SetEditFile(app))
	r.Delete("/file/{id}", handlers.DeleteFile(app))

	r.With(LimitUploadSize(config.MAX_UPLOAD_SIZE)).Post("/upload", handlers.UploadFile(app))
	return r
}
