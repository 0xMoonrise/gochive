package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (db *DBhdlr) GetFile(c *gin.Context) {

	p := c.Param("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong"})
		slog.Warn("The id param cannot convert to int")
		return
	}

	data, err := db.Query.GetArchive(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
		return
	}

    c.Header("Content-Type", "application/pdf")

	c.Writer.WriteHeader(http.StatusOK)
    c.Writer.Write(data.File)
  
}
