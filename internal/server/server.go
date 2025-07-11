package server

import (
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *database.Queries) *gin.Engine {

	r := gin.Default()

	hdlr := handlers.Handler(cfg)

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.Root)
	r.GET("/view", handlers.ViewFile)
	r.GET("/get_file/:id", hdlr.GetFile)
	r.GET("/get_files/:page", hdlr.GetFiles)

	r.POST("/upload", hdlr.UploadFile)
	r.POST("/search/:page", hdlr.SearchFiles)
	r.POST("/set_favorite/:id", hdlr.SetFavorite)
	r.POST("/edit/:id", hdlr.SetEditFile)

	return r
}
