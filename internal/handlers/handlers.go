package handlers

import (
	"net/http"
	"path"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

const pageSize = 8

func GetImage(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		param := c.Param("name")
		objKey := path.Join("images", param)
		length, contentType, reader, err := app.Storage.GetItem(c.Request.Context(), objKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		defer reader.Close()
		c.DataFromReader(http.StatusOK, length, contentType, reader, nil)
	}
}
