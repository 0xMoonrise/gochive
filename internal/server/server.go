package server

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/gin-gonic/gin"
)

var cachedCSS []byte
var once sync.Once

func InjectUserCSS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "viewer.css") {
			once.Do(func() {
				viewerCSS, _ := os.ReadFile("/opt/pdfjs/web/viewer.css")
				userCSS, _ := os.ReadFile("static/styles/userContent.css")
				cachedCSS = append(append(viewerCSS, "\n\n/* userContent.css */\n"...), userCSS...)
			})
			c.Data(http.StatusOK, "text/css", cachedCSS)
			c.Abort()
			return
		}
		c.Next()
	}
}

func NewServer(cfg *database.Queries) *gin.Engine {

	r := gin.Default()

	// change this if you need to trust proxies
	// see https://gin-gonic.com/es/docs/deployment/#dont-trust-all-proxies
	r.SetTrustedProxies(nil)

	hdlr := handlers.Handler(cfg)

	r.LoadHTMLGlob("templates/*")

	r.Use(InjectUserCSS())

	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "static/favicon.ico")
	r.StaticFS("/build", http.Dir("/opt/pdfjs/build"))
	r.StaticFS("/web", http.Dir("/opt/pdfjs/web/"))

	r.GET("/", handlers.Root)
	r.GET("/view", handlers.ViewFile)
	r.GET("/:id", hdlr.GetFile)

	r.GET("/get_files/:page", hdlr.GetFiles)

	r.POST("/upload", hdlr.UploadFile)
	r.POST("/search/:page", hdlr.SearchFiles)
	r.POST("/set_favorite/:id", hdlr.SetFavorite)
	r.POST("/edit/:id", hdlr.SetEditFile)

	return r
}
