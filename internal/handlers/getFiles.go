package handlers

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/database"
)

func (db *DBhdlr) GetFiles(c *gin.Context) {
	var pageSize int32

	p, err := strconv.Atoi(c.Param("page")) // is this the right way? 
	page := int32(p)
	
	if err != nil {
		slog.Warn("Error trying to parse the page number")
		c.JSON(http.StatusInternalServerError, gin.H{"status":"Something went wrong... "})
	}
	
	pageSize = 8 // How many items will be display
	
	pageElements, _ := db.Query.GetCountArchive(c)
	pageLimit := int32(math.Ceil(float64(pageElements) / float64(pageSize)))

	
	if (page <= 0) || (page > int32(pageLimit)) {
		c.JSON(http.StatusNotFound, gin.H{"status":"page not found"})
		return 
	}

	pageDb, err :=  db.Query.GetArchivePage(c, database.GetArchivePageParams{
		Limit: pageSize,
		Offset: (page - 1) * pageSize,
	})

	if err != nil {

		slog.Error("Error fetching the data from database.")
		c.JSON(http.StatusInternalServerError, gin.H{"status":"Something went wrong..."})

		return
	}
	
	c.JSON(http.StatusOK, gin.H{"files":pageDb, "pages":pageLimit})

	
}
