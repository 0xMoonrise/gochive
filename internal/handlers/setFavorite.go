package handlers

import (
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *App) SetFavorite(c *gin.Context) {

	param1, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		slog.Warn("Error trying to parse the page number")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong... "})
		return
	}

	id := int32(param1)
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
