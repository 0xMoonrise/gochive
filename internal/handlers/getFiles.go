package handlers

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func GetFiles(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.ParseInt(c.Param("page"), 10, 64)

		if err != nil {
			slog.Error("Error trying to parse the page number")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong... "})
			return
		}

		pageElements, err := app.Db.GetCountArchive(c)
		if err != nil {
			slog.Error("Something went wrong while trying to fetch data from database", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong... "})
			return
		}
		pageLimit := int64(math.Ceil(float64(pageElements) / float64(pageSize)))

		if (page <= 0) || (page > pageLimit) {
			c.JSON(http.StatusNotFound, gin.H{"status": "page not found"})
			return
		}

		pageDb, err := app.Db.GetArchivePage(c, database.GetArchivePageParams{
			Limit:  pageSize,
			Offset: (page - 1) * pageSize,
		})

		if err != nil {
			slog.Error("Error fetching the data from database.")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong..."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": pageDb, "pages": pageLimit})

	}
}
