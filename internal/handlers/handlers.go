package handlers

import (
	"net/http"
	"path"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/gin-gonic/gin"
)

type App struct {
	Db       *database.Queries
	S3       string
	S3Client *utils.Bucket
}

func (app *App) GetImage(c *gin.Context) {
	param := c.Param("name")
	objKey := path.Join("images", param)

	err := app.S3Client.StreamFile(c.Request.Context(), app.S3, objKey, c.Writer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Header("Content-Type", "image/webp")
}
