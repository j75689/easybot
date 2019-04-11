package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionMiddleware handle dashboard session
func SessionMiddleware(contextPath string) gin.HandlerFunc {

	return func(c *gin.Context) {
		var isLogin = false
		session := sessions.Default(c)
		v := session.Get("login")
		if v != nil {
			isLogin = v.(bool)
		}
		if isLogin || strings.Index(c.Request.URL.Path, "/login") > -1 {
			c.Next()
		} else {
			c.Redirect(http.StatusMovedPermanently, contextPath+"login")
			c.Abort()
		}

	}

}
