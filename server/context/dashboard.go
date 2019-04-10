package context

import (
	"net/http"

	"github.com/gin-contrib/sessions"
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

// HandleLogin process dashboard login
func HandleLogin(adminUser, adminPass string) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			user   = c.DefaultPostForm("user", "")
			pass   = c.DefaultPostForm("pass", "")
			result = gin.H{
				"success": false,
			}
		)
		if user == adminUser && pass == adminPass {
			session := sessions.Default(c)
			session.Set("login", true)
			session.Save()
			result["success"] = true
		}

		c.JSON(200, result)
	}

}
