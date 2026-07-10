package handlers

import (
	"database/sql"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func SearchFiles(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		search := c.PostForm("search")
		page, err := strconv.ParseInt(c.Param("page"), 10, 64)
		if err != nil {
			slog.Error("cannot convert the page parameter on search file", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
			return
		}

		s := sql.NullString{
			String: search,
			Valid:  true,
		}

		pageElements, err := app.DB.Queries.GetCountSearch(c, s)
		if err != nil {
			slog.Error("cannot fetch the count", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wrong..."})
			return
		}

		pageLimit := math.Ceil(float64(pageElements) / float64(config.PAGE_SIZE))
		if (page <= 0) || (page > int64(pageLimit)) {
			c.JSON(http.StatusNotFound, gin.H{"status": "page not found"})
			return
		}

		searchParam := database.SearchArchiveParams{
			Column1: s,
			Limit:   config.PAGE_SIZE,
			Offset:  int(page-1) * config.PAGE_SIZE,
		}

		data, err := app.DB.Queries.SearchArchive(c, searchParam)

		if err != nil {
			slog.Error("cannot fetch the data from database", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "something went wront..."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": data, "pages": pageLimit})
	}
}
