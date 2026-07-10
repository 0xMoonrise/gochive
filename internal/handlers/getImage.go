package handlers

import (
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

func GetImage(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "invalid id"})
			return
		}

		objKey := path.Join("images", id)
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
