package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ViewFile(c *gin.Context) {
    c.HTML(http.StatusOK, "view_pdf.html", gin.H{
      "title": "Main website",
    })
}
