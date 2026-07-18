package handlers

import (
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
)

func GetImage(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil || id <= 0 {
			slog.Warn("invalid id param", "raw", paramID)
			Error(w, http.StatusBadRequest, "Something went wrong")
			return
		}
		if _, err := app.DB.Queries.GetArchiveById(r.Context(), id); err != nil {
			Error(w, http.StatusNotFound, "Not found")
			return
		}

		objKey := path.Join("images", paramID)
		obj, err := app.Storage.GetItem(r.Context(), objKey)
		if err != nil {
			slog.Warn("image not found", "error", err)
			Error(w, http.StatusNotFound, "not found")
			return
		}

		defer obj.Reader.Close()

		if err := fromStorageObject(w, obj); err != nil {
			slog.Error("failed to stream file to response", "error", err, "id", id)
		}

	}
}
