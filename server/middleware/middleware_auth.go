package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/server/auth"
)

// UserAuthMiddleware handle auth
func UserAuthMiddleware(skipper RouteSkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		if skipper(c) {
			c.Next()
			return
		}

		token, err := auth.GetTokenFromRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if token.Subject == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("user not login")})
			c.Abort()
			return
		}

	}
}
