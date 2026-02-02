package handlers

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/database"
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

		fileReader, err := file.Open()
		if err != nil {
			return
		}

		defer fileReader.Close()
		err = app.Storage.PutItem(
			c.Request.Context(),
			path.Join("files", file.Filename),
			file.Size,
			"application/octet-stream",
			fileReader,
		)

		if err != nil {
			slog.Error("Error while trying to upload a file to storage", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Uploaded unsuccessful",
			})
			return
		}

		id, err := app.Db.InsertFile(c, database.InsertFileParams{
			Filename:  file.Filename,
			Editorial: "Default",
		})

		if err != nil {
			slog.Error("Error while trying to store metada file into the database", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Uploaded unsuccessful",
			})
			return
		}

		image := &bytes.Buffer{}
		err = utils.MakeThumbnail(fileReader, file.Size, 0, image)
		if err != nil {
			slog.Error("Error while trying to generate the thumbnail", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Uploaded unsuccessful",
			})
			return
		}

		imageBytes := image.Bytes()
		imageReader := bytes.NewReader(imageBytes)
		size := imageReader.Size()

		objKey := path.Join("images", strconv.Itoa(int(id)))
		err = app.Storage.PutItem(c.Request.Context(), objKey, size, "application/octet-stream", imageReader)

		if err != nil {
			slog.Error("Error while trying to upload a image to storage", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Uploaded unsuccessful",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Uploaded successful"})
	}
}
