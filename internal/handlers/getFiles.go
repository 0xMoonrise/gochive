package handlers

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func GetFiles(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Param("page"))
		if err != nil {
			slog.Error("Error trying to parse the page number")
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong... "})
			return
		}

		pageElements, err := app.DB.Queries.GetCountArchive(c)
		if err != nil {
			slog.Error("Something went wrong while trying to fetch data from database", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong... "})
			return
		}

		pageLimit := int(math.Ceil(float64(pageElements) / float64(config.PAGE_SIZE)))
		if (page <= 0) || (page > pageLimit) {
			c.JSON(http.StatusNotFound, gin.H{"status": "page not found"})
			return
		}

		pageDb, err := app.DB.Queries.GetArchivePage(c,
			database.GetArchivePageParams{
				Limit:  config.PAGE_SIZE,
				Offset: (page - 1) * config.PAGE_SIZE,
			})

		if err != nil {
			slog.Error("Error fetching the data from database.")
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": pageDb, "pages": pageLimit})

	}
}
