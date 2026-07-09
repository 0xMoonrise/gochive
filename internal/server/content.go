package server

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var cachedCSS []byte
var once sync.Once

func injectConentCss() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "viewer.css") {
			once.Do(func() {
				viewerCSS, err := os.ReadFile("/opt/gochive/lib/pdfjs/web/viewer.css")
				if err != nil {
					log.Printf("InjectUserCSS: failed to read viewer.css: %v", err)
					return
				}
				userCSS, err := os.ReadFile("static/styles/userContent.css")
				if err != nil {
					log.Printf("InjectUserCSS: failed to read userContent.css: %v", err)
					return
				}
				cachedCSS = append(append(viewerCSS, "\n\n/* userContent.css */\n"...), userCSS...)
			})

			if len(cachedCSS) == 0 {
				c.Status(http.StatusInternalServerError)
				c.Abort()
				return
			}
			c.Data(http.StatusOK, "text/css", cachedCSS)
			c.Abort()
			return
		}

		c.Next()
	}
}
