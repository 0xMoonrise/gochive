package handlers

import (
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"io"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
)

func makeThumbnail(data []byte, thumb *[]byte, filename string,  c *gin.Context) error {

	thumbnail, err := utils.GenerateWebpThumbnail(data, "static/thumbnails/")

	if err != nil {
		
		slog.Error("cannot generate the thumbnail")
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status":"something went wrong"})

		return err
	}

	thumName := strings.Replace(filename, "pdf", "webp", 1)
	err = utils.SaveThumbnailToStatic(thumbnail, thumName)

	if err != nil {

		slog.Error("cannot generate the thumbnail")
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status":"something went wrong"})

		return err
	}
	
	*thumb = thumbnail

	return nil
}

func validateFilename(filename string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9._ -]+.(pdf|md)$", filename)
	log.Println(match)
	return match
}


func (db *DBhdlr) UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	var thumbnail []byte
	
	if err != nil {
		slog.Error("something went wrong while uploading the a file")
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Something went wrong"})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "File is empty"})
		return
	}

	if ! validateFilename(file.Filename) {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "File type is not allowed"})
		return
	}
	
	rawFile, err := file.Open()

	if err != nil {
		return 
	}

	defer rawFile.Close()
	
	data, err := io.ReadAll(rawFile)

	if err != nil {
		return
	}

	thumbnail = nil

	if strings.Contains(file.Filename, "pdf") {
		err := makeThumbnail(data, &thumbnail, file.Filename, c)
		if err != nil {
			return 
		}
	}
	
	insertFile := database.InsertFileParams{
		Filename: file.Filename,
		Editorial: "Default",
		File: data,
		ThumbnailImage: thumbnail,
	}
	
	err = db.Query.InsertFile(c, insertFile)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "DB insert failed"})
        return
    }
	
	c.JSON(http.StatusOK, gin.H{"status":"Uploaded successful"})
}
