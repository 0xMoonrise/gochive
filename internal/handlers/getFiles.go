package handlers

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
)

func GetFiles(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, err := strconv.Atoi(r.PathValue("page"))
		if err != nil || page <= 0 {
			slog.Error("error trying to parse the page number")
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}

		pageElements, err := app.DB.Queries.GetCountArchive(r.Context())
		if err != nil {
			slog.Error("something went wrong while trying to fetch data from database", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}

		pageLimit := int(math.Ceil(float64(pageElements) / float64(config.PAGE_SIZE)))
		if (page <= 0) || (page > pageLimit) {
			Error(w, http.StatusNotFound, "page not found")
			return
		}

		pageDb, err := app.DB.Queries.GetArchivePage(r.Context(),
			database.GetArchivePageParams{
				Limit:  config.PAGE_SIZE,
				Offset: (page - 1) * config.PAGE_SIZE,
			})

		if err != nil {
			slog.Error("error fetching the data from database", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}

		JSON(w, http.StatusOK, FilesResponse{
			Files: pageDb,
			Pages: pageLimit,
		})

	}
}
