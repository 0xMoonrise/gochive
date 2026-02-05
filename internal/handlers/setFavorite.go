package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
)

func SetFavorite(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			slog.Warn("Error trying to parse the page number")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong... "})
			return
		}

		favorite, err := strconv.ParseBool(c.PostForm("favorite"))

		if err != nil {
			slog.Warn("Error trying to parse the favorite bool")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong..."})
			return
		}

		app.Db.SetFavorite(c, database.SetFavoriteParams{
			Favorite: favorite,
			ID:       id,
		})

		c.JSON(http.StatusOK, gin.H{"status": "Favorite Updated"})
	}
}
