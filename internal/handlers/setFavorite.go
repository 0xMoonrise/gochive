package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/database"
	"strconv"
	"log/slog"
)

func (db *DBhdlr) SetFavorite(c *gin.Context) {
	
	param1, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		slog.Warn("Error trying to parse the page number")
		c.JSON(http.StatusInternalServerError, gin.H{"status":"Something went wrong... "})
		return
	}

	id := int32(param1)
	favorite, err := strconv.ParseBool(c.PostForm("favorite"))

	if err != nil {
		slog.Warn("Error trying to parse the favorite bool")
		c.JSON(http.StatusInternalServerError, gin.H{"status":"Something went wrong..."})
		return
	}
	
	db.Query.SetFavorite(c, database.SetFavoriteParams{
		Favorite: favorite,
		ID: id,
	})

	c.JSON(http.StatusOK, gin.H{"status":"Favorite Updated"})
}
