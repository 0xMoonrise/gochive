package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func SetFavorite(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("error trying to parse the page number", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong... "})
			return
		}

		favorite, err := strconv.ParseBool(c.PostForm("favorite"))
		if err != nil {
			slog.Warn("error trying to parse the favorite bool", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
			return
		}

		err = app.DB.Queries.SetFavorite(c,
			database.SetFavoriteParams{
				Favorite: favorite,
				ID:       id,
			})

		if err != nil {
			slog.Error("error trying to set the value to favorite", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Favorite Updated"})
	}
}
