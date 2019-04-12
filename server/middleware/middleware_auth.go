package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/auth"
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

		token, claim, err := auth.GetTokenFromRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("claim", claim)

	}
}
