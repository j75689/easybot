package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/auth/claim"
	"github.com/j75689/easybot/auth/token"
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/model"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

// ScopeMiddleware handle valid api scope
func ScopeMiddleware(db *store.Storage, skipper RouteSkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		if skipper != nil {
			if skipper(c) {
				c.Next()
				return
			}
		}

		var (
			claims    *claim.ServiceAccountClaims
			tokenInfo *token.TokenInfo
		)
		claimsObj, claimsOk := c.Get("claim")
		tokenObj, tokenOk := c.Get("token")

		if claimsOk && tokenOk {
			tokenInfo = tokenObj.(*token.TokenInfo)
			claims = claimsObj.(*claim.ServiceAccountClaims)
			value, _ := (*db).Load(config.ServiceAccountTable, claims.Name)
			var account model.ServiceAccount
			if data, err := json.Marshal(value); err == nil {
				json.Unmarshal(data, &account)
				// Verify Account Info
				if account.ValidInfo(tokenInfo, claims) {
					// Verify Scope
					path := c.Request.URL.Path
					if !config.Scope.Allow(claims.Scope, path) {
						c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
						c.Abort()
					}
				} else {
					logger.Errorf("[scope] Token info not match account [%v]", account)
					c.JSON(http.StatusForbidden, gin.H{"error": "Token info invalid"})
					c.Abort()
				}
			} else {
				logger.Errorf("[scope] unmarshal account [%v] error [%v]", claims.Name, err)
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Token info not found"})
			c.Abort()
		}
	}
}
