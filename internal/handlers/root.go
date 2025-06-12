package handlers

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
      "title": "Main website",
    })
}
