package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
)

func SetFavorite(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Warn("error trying to parse the page number", "error", err)
			Error(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		favorite, err := strconv.ParseBool(r.FormValue("favorite"))
		if err != nil {
			slog.Warn("error trying to parse the favorite bool", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}

		err = app.DB.Queries.SetFavorite(r.Context(),
			database.SetFavoriteParams{
				Favorite: favorite,
				ID:       id,
			})

		if err != nil {
			slog.Error("error trying to set the value to favorite", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}

		JSON(w, http.StatusOK, Success{
			Status: "favorite updated",
		})
	}
}
