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

func View(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			slog.Error("cannot convert id on view", "error", err)
			Error(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		filename, err := app.DB.Queries.GetArchiveById(r.Context(), id)
		if err != nil {
			slog.Error("id not found on view", "error", err)
			Error(w, http.StatusNotFound, "Not found")
			return
		}

		if strings.HasSuffix(filename, ".md") {
			render(w, app.Templates, "view_md.html", http.StatusOK, map[string]any{
				"title": filename,
				"id":    id,
			})
			return
		}

		vendor, err := loadVendorViewer()
		if err != nil {
			slog.Error("load vendor failed on view", "error", err)
			Error(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		inject, err := renderHeadInject(filename, "/file/"+paramID)
		if err != nil {
			slog.Error("cannot inject content on view", "error", err)
			Error(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		html := bytes.Replace(vendor, []byte("<head>"), inject, 1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(html); err != nil {
			slog.Error("failed to stream file to response", "error", err, "id", id)
		}

	}
}
