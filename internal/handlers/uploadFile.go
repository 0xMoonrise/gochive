package handlers

import (
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-sanitize"
)

func (db *DBhdlr) UploadFile(c *gin.Context) {

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

	rawData, err := io.ReadAll(rawFile)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
		return
	}

	filename := sanitize.XSS(file.Filename)

	insertFile := database.InsertFileParams{
		Filename:  filename,
		Editorial: "Default",
		File:      rawData,
	}

	id, err := db.Query.InsertFile(c, insertFile)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "DB insert failed"})
		return
	}

	if strings.Contains(file.Filename, "pdf") {
		thumbnail, err := utils.MakeThumbnail(rawData, strconv.Itoa(int(id)))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
			return
		}

		err = db.Query.SaveThumbnail(c, database.SaveThumbnailParams{
			ID:             id,
			ThumbnailImage: thumbnail,
		})

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong..."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Uploaded successful"})
}
