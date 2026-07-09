package handlers

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

func GetImage(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		paramName := c.Param("name")
		objKey := path.Join("images", paramName)
		obj, err := app.Storage.GetItem(c.Request.Context(), objKey)
		if err != nil {
			slog.Warn("Image not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
			return
		}

		defer obj.Reader.Close()
		c.DataFromReader(http.StatusOK, obj.Length, obj.ContentType, obj.Reader, nil)
	}
}
