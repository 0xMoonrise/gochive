package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (db *DBhdlr) GetFile(c *gin.Context) {

	p := c.Param("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"}) // check status request
		slog.Warn("The id param cannot convert to int")
		return
	}

	data, err := db.Query.GetArchive(c, int32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
		return
	}
	// Maybe support other kind of files?
	if strings.Contains(data.Filename, "pdf") {

		c.Header("Content-Type", "application/pdf")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(data.File)

	}

	if strings.Contains(data.Filename, "md") { // should I assume that I will never store anything else but md and pdf?

		html := markdown.ToHTML(data.File, nil, nil)

		c.HTML(http.StatusOK, "view_md.html", gin.H{
			"Content": template.HTML(html),
		})

	}

}
