package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/server/auth"
)

// UserAuthMiddleware handle auth
func UserAuthMiddleware(skipper RouteSkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		if skipper != nil {
			if skipper(c) {
				c.Next()
				return
			}
		}

		token, err := auth.GetTokenFromRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		// verify scope
		_ = token

	}
}
