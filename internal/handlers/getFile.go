package handlers

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

func GetFile(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Something went wrong"}) // check status request
			slog.Warn("The id param cannot convert to int")
			return
		}

		filename, err := app.Db.GetArchiveById(c, id)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
			return
		}

		objKey := path.Join("files", filename)
		length, contentType, reader, err := app.Storage.GetItem(c.Request.Context(), objKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
		}

		defer reader.Close()
		if strings.Contains(filename, "md") {
			data, _ := io.ReadAll(reader)
			html := markdown.ToHTML(data, nil, nil)

			c.HTML(http.StatusOK, "view_md.html", gin.H{
				"Content": template.HTML(html),
			})
			return
		}

		c.DataFromReader(http.StatusOK, length, contentType, reader, nil)
	}
}
