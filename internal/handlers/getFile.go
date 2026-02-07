package handlers

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

func GetFile(app *core.App) gin.HandlerFunc {
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
		obj, err := app.Storage.GetItem(c.Request.Context(), objKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
		}

		defer obj.Reader.Close()
		if strings.Contains(filename, "md") {
			data, _ := io.ReadAll(obj.Reader)
			html := markdown.ToHTML(data, nil, nil)

			c.HTML(http.StatusOK, "view_md.html", gin.H{
				"Content": template.HTML(html),
			})
			return
		}

		c.DataFromReader(http.StatusOK, obj.Length, obj.ContentType, obj.Reader, nil)
	}
}
