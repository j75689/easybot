package context

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleIndexPage dashboard index page.
func HandleIndexPage(context_path string) func(*gin.Context) {

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Redirect(301, context_path+"dashboard")
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}
