package handlers

import (
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"io"
	"github.com/gin-gonic/gin"
	"github.com/0xMoonrise/gochive/internal/database"
)

func validateFilename(filename string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9._ -]+.(pdf|md)$", filename)
	log.Println(match)
	return match
}


func (db *DBhdlr) UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")

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
	
	insertFile := database.InsertFileParams{
		Filename: file.Filename,
		Editorial: "Default",
		File: data,
	}
	
	err = db.Query.InsertFile(c, insertFile)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "DB insert failed"})
        return
    }
	
	c.JSON(http.StatusOK, gin.H{"status":"Uploaded successful"})
}
