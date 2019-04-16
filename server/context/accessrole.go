package context

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/auth"
	"github.com/j75689/easybot/auth/token"
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/model"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

// HandleGetAllServiceAccount process get all service account info
func HandleGetAllServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var accounts = []model.ServiceAccount{}
		err := (*db).LoadAll(config.ServiceAccountTable, func(id string, value interface{}) {
			var account model.ServiceAccount
			if data, err := json.Marshal(value); err == nil {
				json.Unmarshal(data, &account)
				accounts = append(accounts, account)
			} else {
				logger.Errorf("[dashboard] unmarshal account [%v] error [%v]", id, err)
			}
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, accounts)

	}
}

// HandleGetServiceAccount process get service account info
func HandleGetServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			account model.ServiceAccount
			name    = c.Param("name")
		)
		value, err := (*db).LoadWithFilter(config.ServiceAccountTable, map[string]interface{}{"name": name})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Account not found"})
			return
		}
		if data, err := json.Marshal(value); err == nil {
			json.Unmarshal(data, &account)
			c.JSON(http.StatusOK, account)
		} else {
			logger.Errorf("[dashboard] unmarshal account [%v] error [%v]", name, err)
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Data error"})
		}
	}
}

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
		// Check Exist
		if _, err := (*db).LoadWithFilter(config.ServiceAccountTable, map[string]interface{}{"name": name}); err == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": name + " is Existed"})
			return
		}

		// Create New
		active, _ = strconv.Atoi(activeStr)

		account := model.ServiceAccount{
			Name:     name,
			EMail:    email,
			Domain:   domain,
			Provider: provider,
			Active:   int64(active),
			Scope:    scope,
		}
		token, err := auth.GenerateToken(&account)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		account.Generate = time.Now().Unix()
		account.Token = token.AccessToken
		account.Expired = token.Expire

		if err := (*db).Save(config.ServiceAccountTable, &account); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

// HandleSaveServiceAccount process save service account
func HandleSaveServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			account   model.ServiceAccount
			name      = c.Param("name")
			newName   = c.DefaultPostForm("name", name)
			email     = c.DefaultPostForm("email", "")
			domain    = c.DefaultPostForm("domain", "")
			provider  = c.DefaultPostForm("provider", "")
			activeStr = c.DefaultPostForm("active", "7200")
			active    int
			scope     = c.DefaultPostForm("scope", "")
		)

		value, err := (*db).LoadWithFilter(config.ServiceAccountTable, map[string]interface{}{"name": name})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Account not found"})
			return
		}
		if data, err := json.Marshal(value); err == nil {
			active, _ = strconv.Atoi(activeStr)
			json.Unmarshal(data, &account)
			account.Name = newName
			account.EMail = email
			account.Domain = domain
			account.Provider = provider
			account.Active = int64(active)
			account.Scope = scope

			token, err := auth.GenerateToken(&account)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
				return
			}
			account.Token = token.AccessToken
			account.Expired = token.Expire

			// nmae changed
			if name != newName {
				(*db).Delete(config.ServiceAccountTable, account.ID)
			}

			if err := (*db).SaveWithFilter(config.ServiceAccountTable, &account, map[string]interface{}{"name": account.Name}); err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})
		} else {
			logger.Errorf("[dashboard] unmarshal account [%v] error [%v]", name, err)
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Data error"})
		}
	}
}

// HandleBatchDeleteServiceAccount process batch delete service account
func HandleBatchDeleteServiceAccount(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			accounts []string
		)
		formdata, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(formdata, &accounts)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}

		for _, account := range accounts {
			err = (*db).DeleteWithFilter(config.ServiceAccountTable, map[string]interface{}{"name": account})
			if err != nil {
				logger.Error("[dashboard] delete account [%s] error [%v]", account, err.Error())
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})

	}
}

// HandleRefreshServiceAccountToken process refresh token
func HandleRefreshServiceAccountToken(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			name      = c.Param("name")
			tokenInfo *token.TokenInfo
			account   model.ServiceAccount
		)
		value, _ := (*db).LoadWithFilter(config.ServiceAccountTable, map[string]interface{}{"name": name})

		if data, err := json.Marshal(value); err == nil {
			json.Unmarshal(data, &account)
			tokenInfo, err = auth.GenerateToken(&account)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
				return
			}
			logger.Info(tokenInfo)
			// Save
			account.Token = tokenInfo.AccessToken
			account.Expired = tokenInfo.Expire
			if err := (*db).SaveWithFilter(config.ServiceAccountTable, &account, map[string]interface{}{"name": account.Name}); err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"token":   tokenInfo,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Account not found"})
		}

	}
}

// HandleGetScopeTags process get scope tags
func HandleGetScopeTags() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, config.Scope.Tags())
	}
}
