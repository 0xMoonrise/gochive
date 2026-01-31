package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

func GetFile(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Param("id")
		id, err := strconv.Atoi(p)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"}) // check status request
			slog.Warn("The id param cannot convert to int")
			return
		}

		data, err := app.Db.GetArchive(c, int32(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
			return
		}
		// Maybe support other kind of files?
		if strings.Contains(data.Filename, "pdf") {
			c.Data(http.StatusOK, "application/pdf", data.File)
		}

		if strings.Contains(data.Filename, "md") { // should I assume that I will never store anything else but md and pdf?

			html := markdown.ToHTML(data.File, nil, nil)

			c.HTML(http.StatusOK, "view_md.html", gin.H{
				"Content": template.HTML(html),
			})

		}

	}
}
