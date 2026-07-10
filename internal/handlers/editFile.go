package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-sanitize"
)

func SetEditFile(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.PostForm("filename")
		editorial := c.PostForm("editorial")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Error("cannot convert the id parameter", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
			return
		}

		if !utils.ValidateFilename(filename) {
			slog.Warn("invalid filename", "id", id)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Extension not allowed"})
			return
		}

		filename = sanitize.XSS(filename)
		editorial = sanitize.XSS(editorial)

		err = app.DB.Queries.SetEditFile(c, database.SetEditFileParams{
			Filename:  filename,
			Editorial: editorial,
			ID:        id,
		})

		if err != nil {
			slog.Error("cannot update the values in to the data base", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Resource updated"})
	}
}
