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
)

func SearchFiles(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		search := r.FormValue("search")
		page, err := strconv.Atoi(r.PathValue("page"))
		if err != nil {
			slog.Error("cannot convert the page parameter on search file", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong...")
			return
		}

		s := sql.NullString{
			String: search,
			Valid:  true,
		}

		pageElements, err := app.DB.Queries.GetCountSearch(r.Context(), s)
		if err != nil {
			slog.Error("cannot fetch the count", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong...")
			return
		}

		pageLimit := int(math.Ceil(float64(pageElements) / float64(config.PAGE_SIZE)))
		if (page <= 0) || (page > pageLimit) {
			Error(w, http.StatusNotFound, "page not found")
			return
		}

		searchParam := database.SearchArchiveParams{
			Column1: s,
			Limit:   config.PAGE_SIZE,
			Offset:  int(page-1) * config.PAGE_SIZE,
		}

		searched, err := app.DB.Queries.SearchArchive(r.Context(), searchParam)

		if err != nil {
			slog.Error("cannot fetch the data from database", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong...")
			return
		}

		JSON(w, http.StatusOK, FilesResponse{
			Files: searched,
			Pages: pageLimit,
		})
	}
}
