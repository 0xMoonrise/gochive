package server

import (
    "github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/handlers"
)

func NewServer() *gin.Engine {

	r := gin.Default()

	r.GET("/hello", handlers.Hello)
 
    return r
}
