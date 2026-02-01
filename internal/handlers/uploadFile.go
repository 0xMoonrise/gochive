package handlers

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/gin-gonic/gin"
)

func UploadFile(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			slog.Error("something went wrong while uploading the a file")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"})
			return
		}

		if file.Size == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "File is empty"})
			return
		}

		if !utils.ValidateFilename(file.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "File type is not allowed"})
			return
		}

		rawFile, err := file.Open()

		if err != nil {
			return
		}

		defer rawFile.Close()

		c.JSON(http.StatusOK, gin.H{"status": "Uploaded successful"})
	}
}
