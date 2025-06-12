package server

import (
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *database.Queries) *gin.Engine {
	
	r := gin.Default()

	hdlr := handlers.Handler(cfg)
	
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.Root)
	r.GET("/ping", handlers.Ping)
	r.GET("/get_name/:id", hdlr.GetName)

	return r
}
