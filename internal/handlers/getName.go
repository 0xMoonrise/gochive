package handlers

import (
	"database/sql"
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

	param := sql.NullInt32{Int32: int32(id), Valid: true }
	name, err := db.Query.GetArchive(c, param)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"name": name})
}
