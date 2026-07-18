package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
)

func SetEditFile(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("filename")
		editorial := r.FormValue("editorial")

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("cannot convert the id parameter", "error", err)
			Error(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		if !utils.ValidateFilename(filename) {
			slog.Warn("invalid filename", "id", id)
			Error(w, http.StatusBadRequest, "extension not allowed")
			return
		}

		if utils.IsTooLong(editorial) {
			Error(w, http.StatusBadRequest, "editorial too long")
			return
		}

		if utils.IsTooLong(filename) {
			Error(w, http.StatusBadRequest, "filename too long")
			return
		}

		err = app.DB.Queries.SetEditFile(r.Context(), database.SetEditFileParams{
			Filename:  filename,
			Editorial: editorial,
			ID:        id,
		})

		if err != nil {
			slog.Error("cannot update the values in to the data base", "error", err)
			Error(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		JSON(w, http.StatusOK, Success{
			Status: "resource updated",
		})
	}
}
