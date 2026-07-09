package handlers

import (
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

func DeleteFile(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Warn("Error trying to parse the page number",
				"error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong... "})
			return
		}

		filename, err := app.DB.Queries.GetArchiveById(c, id)
		if err != nil {
			slog.Error("Something went wrong fetching the file on DeleteFile",
				"error", err)
			c.JSON(http.StatusNotModified, gin.H{"status": "something went wrong..."})
			return
		}

		if err := app.DB.Queries.DeleteFile(c, id); err != nil {
			slog.Error("Something went wrong while trying to delete a file",
				"error", err)
			c.JSON(http.StatusNotModified, gin.H{"status": "something went wrong..."})
			return
		}

		objKey := path.Join("files", idParam)
		if err := app.Storage.DelItem(c, objKey); err != nil {
			slog.Error("Something went wrong while trying to delete file from storage",
				"error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
			return
		}

		if !strings.HasSuffix(filename, ".md") {
			objKey = path.Join("images", idParam)
			if err := app.Storage.DelItem(c, objKey); err != nil {
				slog.Error("Something went wrong while trying to delete image from storage",
					"error", err)
				c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "Deleted successfuly"})
	}
}
