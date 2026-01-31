package handlers

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/gin-gonic/gin"
)

func GetImage(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		param := c.Param("name")
		objKey := path.Join("images", param)
		bucket := os.Getenv("BUCKET")

		c.Header("Content-Type", "image/webp")
		err := app.S3Client.StreamFile(c.Request.Context(), bucket, objKey, c.Writer)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
	}
}
