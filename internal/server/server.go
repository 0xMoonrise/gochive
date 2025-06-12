package server

import (
    "github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/handlers"
)

func NewServer() *gin.Engine {

	r := gin.Default()
	
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.Root)
	r.GET("/hello", handlers.Hello)
	
    return r
}
