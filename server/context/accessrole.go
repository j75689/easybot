package context

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/model"
	"github.com/j75689/easybot/pkg/store"
	"github.com/j75689/easybot/server/auth"
)

// HandleCreateServiceAccount process create service account
func HandleCreateServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			name      = c.Param("name")
			email     = c.DefaultPostForm("email", "")
			domain    = c.DefaultPostForm("domain", "")
			provider  = c.DefaultPostForm("provider", "")
			activeStr = c.DefaultPostForm("active", "7200")
			active    int
			scope     = c.DefaultPostForm("scope", "")
		)

		active, _ = strconv.Atoi(activeStr)

		account := &model.ServiceAccount{
			Name:     name,
			EMail:    email,
			Domain:   domain,
			Provider: provider,
			Active:   active,
			Scope:    scope,
		}
		token, err := auth.GenerateToken(account)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err})
			return
		}
		account.Generate = time.Now().Unix()
		account.Token = token.AccessToken
		account.Expired = token.Expire

		if err := (*db).Save(config.ServiceAccountTable, name, account); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

// HandleDeleteServiceAccount process delete service account
func HandleDeleteServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {

	}
}
