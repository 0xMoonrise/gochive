package handlers

import (
	"database/sql"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func SearchFiles(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		search := c.PostForm("search")
		page, err := strconv.ParseInt(c.Param("page"), 10, 64)
		if err != nil {
			slog.Error("cannot convert the page parameter on search file")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wrong..."})
			return
		}

		s := sql.NullString{
			String: search,
			Valid:  true,
		}

		pageElements, _ := app.Db.GetCountSearch(c, s)
		log.Print(pageElements)
		pageLimit := math.Ceil(float64(pageElements) / float64(pageSize))

		if (page <= 0) || (page > int64(pageLimit)) {
			c.JSON(http.StatusNotFound, gin.H{"status": "page not found"})
			return
		}

		searchParam := database.SearchArchiveParams{
			Column1: s,
			Limit:   pageSize,
			Offset:  (page - 1) * pageSize,
		}

		data, err := app.Db.SearchArchive(c, searchParam)

		if err != nil {
			slog.Error("cannot fetch the data from database")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wront..."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": data, "pages": pageLimit})
	}
}
