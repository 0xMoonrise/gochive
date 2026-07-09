package server

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

func limitUploadSize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func NewEngine() *gin.Engine {
	r := gin.Default()

	// change this if you need to trust proxies
	// see https://gin-gonic.com/es/docs/deployment/#dont-trust-all-proxies
	r.SetTrustedProxies(nil)
	templates := template.Must(template.ParseFS(templatesFS, "templates/*"))
	r.SetHTMLTemplate(templates)

	return r
}

func NewServer(app *core.App) *gin.Engine {
	r := NewEngine()

	r.Use(injectConentCss())
	r.Static("/static", "./static")
	r.Static("/lib", "/opt/gochive/lib")
	r.StaticFile("/favicon.ico", "static/favicon.ico")
	r.StaticFS("/build", http.Dir("/opt/gochive/lib/pdfjs/build/"))
	r.StaticFS("/web", http.Dir("/opt/gochive/lib/pdfjs/web/"))

	r.GET("/", handlers.Root)
	r.GET("/file/:id", handlers.GetFile(app))
	r.GET("/view/:id", handlers.View(app))

	r.GET("/images/:name", handlers.GetImage(app))
	r.GET("/get_files/:page", handlers.GetFiles(app))

	r.POST("/upload",
		limitUploadSize(config.MAX_UPLOAD_SIZE),
		handlers.UploadFile(app))
	r.POST("/search/:page", handlers.SearchFiles(app))
	r.POST("/set_favorite/:id", handlers.SetFavorite(app))

	r.PATCH("/edit/:id", handlers.SetEditFile(app))
	r.DELETE("/file/:id", handlers.DeleteFile(app))
	return r
}
