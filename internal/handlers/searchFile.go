package handlers

import (
	"database/sql"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"
)

func (db *DBhdlr) SearchFiles(c *gin.Context) {

	var pageSize int32

	search := c.PostForm("search")
	pageSize = 8

	p, err := strconv.Atoi(c.Param("page"))

	if err != nil {
		slog.Error("cannot convert the page parameter on search file")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wrong..."})
		return
	}

	s := sql.NullString{
		String: search,
		Valid:  true,
	}

	page := int32(p)

	pageElements, _ := db.Query.GetCountSearch(c, s)
	log.Print(pageElements)
	pageLimit := math.Ceil(float64(pageElements) / float64(pageSize))

	if (page <= 0) || (page > int32(pageLimit)) {
		c.JSON(http.StatusNotFound, gin.H{"status": "page not found"})
		return
	}

	searchParam := database.SearchArchiveParams{
		Column1: s,
		Limit:   pageSize,
		Offset:  (page - 1) * pageSize,
	}

	data, err := db.Query.SearchArchive(c, searchParam)

	if err != nil {
		slog.Error("cannot fetch the data from database")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wront..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": data, "pages": pageLimit})
}
