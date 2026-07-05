package handlers

import (
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

func GetFile(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ParamId := c.Param("id")
		id, err := strconv.Atoi(ParamId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"}) // check status request
			slog.Warn("The id param cannot convert to int")
			return
		}

		if _, err := app.DB.Queries.GetArchiveById(c, id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
			return
		}

		objKey := path.Join("files", ParamId)
		obj, err := app.Storage.GetItem(c.Request.Context(), objKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		defer obj.Reader.Close()
		c.DataFromReader(http.StatusOK, obj.Length, obj.ContentType, obj.Reader, nil)
	}
}
