package handlers

import (
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
)

func DeleteFile(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Warn("error trying to parse the page number",
				"error", err)
			Error(w, http.StatusInternalServerError, "Something went wrong... ")
			return
		}

		filename, err := app.DB.Queries.GetArchiveById(r.Context(), id)
		if err != nil {
			slog.Error("something went wrong fetching the file on DeleteFile",
				"error", err)
			Error(w, http.StatusBadRequest, "Something went wrong... ")
			return
		}

		if err := app.DB.Queries.DeleteFile(r.Context(), id); err != nil {
			slog.Error("something went wrong while trying to delete a file",
				"error", err)
			Error(w, http.StatusBadRequest, "Something went wrong... ")
			return
		}

		objKey := path.Join("files", idParam)
		if err := app.Storage.DelItem(r.Context(), objKey); err != nil {
			slog.Error("something went wrong while trying to delete file from storage",
				"error", err)
			Error(w, http.StatusBadRequest, "Something went wrong... ")
			return
		}

		if !strings.HasSuffix(filename, ".md") {
			objKey = path.Join("images", idParam)
			if err := app.Storage.DelItem(r.Context(), objKey); err != nil {
				slog.Error("something went wrong while trying to delete image from storage",
					"error", err)
				Error(w, http.StatusBadRequest, "Something went wrong... ")
				return
			}
		}

		JSON(w, http.StatusOK, Success{
			Status: "Deleted successfuly",
		})
	}
}
