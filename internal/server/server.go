package server

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/handlers"
)

//go:embed templates/*
var templatesFS embed.FS
var Templates = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

type Middleware func(http.Handler) http.Handler

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", clientIP(r),
			"time", time.Since(start),
		)
	})
}

func NewEngine() *http.ServeMux {
	r := http.NewServeMux()
	return r
}

func NewServer(app *core.App) http.Handler {

	app.Templates = Templates
	loadViewerCSS()
	r := NewEngine()

	fromFS(r, "/static/", "./static")
	fromFS(r, "/lib/", "/opt/gochive/lib")
	fromFS(r, "/build/", "/opt/gochive/lib/pdfjs/build/")
	fromFS(r, "/web/", "/opt/gochive/lib/pdfjs/web/")

	r.HandleFunc("GET /{$}", handlers.Root(app))
	r.HandleFunc("GET /file/{id}", handlers.GetFile(app))
	r.HandleFunc("GET /view/{id}", handlers.View(app))
	r.HandleFunc("GET /images/{id}", handlers.GetImage(app))
	r.HandleFunc("GET /get_files/{page}", handlers.GetFiles(app))

	r.Handle("POST /upload", Chain(
		handlers.UploadFile(app),
		LimitUploadSize(config.MAX_UPLOAD_SIZE),
	))

	r.HandleFunc("POST /search/{page}", handlers.SearchFiles(app))
	r.HandleFunc("POST /set_favorite/{id}", handlers.SetFavorite(app))
	r.HandleFunc("PATCH /edit/{id}", handlers.SetEditFile(app))
	r.HandleFunc("DELETE /file/{id}", handlers.DeleteFile(app))

	r.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})

	return Chain(r, injectContentCSS, loggingMiddleware)
}
