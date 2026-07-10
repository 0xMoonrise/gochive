package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/gin-gonic/gin"
)

var (
	viewerVendor []byte
	viewerOnce   sync.Once
	viewerErr    error
)

type headInjectData struct {
	Title  string
	PDFURL string
}

func loadVendorViewer() ([]byte, error) {
	viewerOnce.Do(func() {
		viewerVendor, viewerErr = os.ReadFile(config.VIEWER_PATH)
	})
	return viewerVendor, viewerErr
}

var headInject = template.Must(template.New("viewer-head-inject").Parse(`<head>
<base href="/web/">
<title>{{.Title}}</title>
<script>
  document.addEventListener("webviewerloaded", function () {
    PDFViewerApplicationOptions.set("defaultUrl", {{.PDFURL}});
    PDFViewerApplication.setTitle = function () {
      document.title = {{.Title}};
    };
  });
</script>`))

func renderHeadInject(title, pdfURL string) ([]byte, error) {
	var buf bytes.Buffer
	if err := headInject.Execute(&buf, headInjectData{Title: title, PDFURL: pdfURL}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func View(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ParamId := c.Param("id")
		id, err := strconv.Atoi(ParamId)
		if err != nil {
			slog.Error("cannot convert id on view", "error", err)
			c.JSON(http.StatusBadRequest, "something went wrong")
			return
		}

		filename, err := app.DB.Queries.GetArchiveById(c, id)
		if err != nil {
			slog.Error("id not found on view", "error", err)
			c.JSON(http.StatusBadRequest, "something went wrong")
			return
		}

		if strings.HasSuffix(filename, ".md") {
			c.HTML(http.StatusOK, "view_md.html", gin.H{
				"title": filename,
				"id":    id,
			})
			return
		}

		vendor, err := loadVendorViewer()
		if err != nil {
			slog.Error("load vendor failed on view", "error", err)
			c.JSON(http.StatusBadRequest, "something went wrong")
			return
		}

		inject, err := renderHeadInject(filename, "/file/"+ParamId)
		if err != nil {
			slog.Error("cannot inject content on view", "error", err)
			c.JSON(http.StatusBadRequest, "something went wrong")
			return
		}

		html := bytes.Replace(vendor, []byte("<head>"), inject, 1)
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	}
}
