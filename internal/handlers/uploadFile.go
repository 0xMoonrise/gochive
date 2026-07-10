package handlers

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/gin-gonic/gin"
)

func UploadFile(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			slog.Error("something went wrong while uploading the a file", "error", err)
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

		// Just in case
		file.Filename = filepath.Base(file.Filename)
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"})
			return
		}
		defer fileReader.Close()

		tx, err := app.DB.Begin()
		if err != nil {
			slog.Error("failed to begin transaction", "error", err)
			return
		}

		defer tx.Rollback()
		qtx := app.DB.Queries.WithTx(tx)
		id, err := qtx.InsertFile(c,
			database.InsertFileParams{
				Filename:  file.Filename,
				Editorial: "Default",
			})

		if err != nil {
			slog.Error("error while trying to store metada file into the database", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Uploaded unsuccessful",
			})
			return
		}

		buffer := make([]byte, 512)
		if int64(len(buffer)) > file.Size {
			buffer = buffer[:file.Size]
		}

		n, err := fileReader.ReadAt(buffer, 0)
		if err != nil && err != io.EOF {
			slog.Error("error reading file for content-type detection", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"})
			return
		}

		contentType := http.DetectContentType(buffer[:n])

		err = app.Storage.PutItem(
			c.Request.Context(),
			path.Join("files", strconv.Itoa(id)),
			&core.Object{
				Length:      file.Size,
				ContentType: contentType,
				Reader:      fileReader,
			},
		)

		if err != nil {
			slog.Error("error while trying to upload a file to storage", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Uploaded unsuccessful",
			})
			return
		}

		if strings.HasSuffix(file.Filename, ".md") {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"file": gin.H{
					"id":        id,
					"filename":  file.Filename,
					"editorial": "Default",
					"favorite":  false,
				},
			})
			tx.Commit()
			return
		}

		image := &bytes.Buffer{}
		err = utils.MakeThumbnail(fileReader, file.Size, 0, image)
		if err != nil {
			slog.Error("error while trying to generate the thumbnail", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Uploaded unsuccessful",
			})
			return
		}

		imageBytes := image.Bytes()
		imageReader := bytes.NewReader(imageBytes)
		size := imageReader.Size()

		sniff := min(len(imageBytes), 512)
		contentType = http.DetectContentType(imageBytes[:sniff])

		objKey := path.Join("images", strconv.Itoa(int(id)))
		err = app.Storage.PutItem(c.Request.Context(), objKey, &core.Object{
			Length:      size,
			ContentType: contentType,
			Reader:      io.NopCloser(imageReader),
		})

		if err != nil {
			slog.Error("error while trying to upload a image to storage", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Uploaded unsuccessful",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"file": gin.H{
				"id":        id,
				"filename":  file.Filename,
				"editorial": "Default",
				"favorite":  false,
			},
		})

		tx.Commit()
	}
}
