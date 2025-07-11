package handlers

import (
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-sanitize"
	"log/slog"
	"net/http"
	"strconv"
)

func (db *DBhdlr) SetEditFile(c *gin.Context) {

	filename := c.PostForm("filename")
	editorial := c.PostForm("editorial")

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		slog.Error("cannot convert the page parameter on search file")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wrong..."})
		return
	}

	if !validateFilename(filename) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Extension not allowed"})
		return
	}

	filename = sanitize.XSS(filename)
	editorial = sanitize.XSS(editorial)

	err = db.Query.SetEditFile(c, database.SetEditFileParams{
		Filename:  filename,
		Editorial: editorial,
		ID:        int32(id),
	})

	if err != nil {
		slog.Error("cannot convert the page parameter on search file")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wrong..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Resource updated"})
}
