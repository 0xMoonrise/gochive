package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (db *DBhdlr) GetName(c *gin.Context) {

	p := c.Param("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong"})
		slog.Warn("The id param cannot convert to int")
		return
	}

	name, err := db.Query.GetArchive(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"name": name})
}
