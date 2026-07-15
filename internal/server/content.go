package server

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var cachedCSS []byte

func loadViewerCSS() error {
	viewerCSS, err := os.ReadFile("/opt/gochive/lib/pdfjs/web/viewer.css")
	if err != nil {
		return errors.New("cannot read viewer.css")
	}

	userCSS, err := os.ReadFile("static/styles/userContent.css")
	if err != nil {
		return errors.New("cannot read userContent.css")
	}

	buf := make([]byte, 0, len(viewerCSS)+len(userCSS)+32)
	buf = append(buf, viewerCSS...)
	buf = append(buf, "\n\n/* userContent.css */\n"...)
	buf = append(buf, userCSS...)

	cachedCSS = buf
	return nil
}

func injectContentCSS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasSuffix(c.Request.URL.Path, "viewer.css") {
			c.Next()
			return
		}
		c.Data(http.StatusOK, "text/css", cachedCSS)
		c.Abort()
	}
}
