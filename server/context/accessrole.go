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
		err := (*db).LoadAll(config.ServiceAccountTable, func(key string, value interface{}) {
			var account model.ServiceAccount
			if data, err := json.Marshal(value); err == nil {
				json.Unmarshal(data, &account)
				accounts = append(accounts, account)
			} else {
				logger.Errorf("[dashboard] unmarshal account [%v] error [%v]", key, err)
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
		value, err := (*db).Load(config.ServiceAccountTable, name)
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
		if _, err := (*db).Load(config.ServiceAccountTable, name); err == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": name + " is Existed"})
			return
		}

		// Create New
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
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		account.Generate = time.Now().Unix()
		account.Token = token.AccessToken
		account.Expired = token.Expire

		if err := (*db).Save(config.ServiceAccountTable, name, account); err != nil {
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
			email     = c.DefaultPostForm("email", "")
			domain    = c.DefaultPostForm("domain", "")
			provider  = c.DefaultPostForm("provider", "")
			activeStr = c.DefaultPostForm("active", "7200")
			active    int
			scope     = c.DefaultPostForm("scope", "")
		)

		value, err := (*db).Load(config.ServiceAccountTable, name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Account not found"})
			return
		}
		if data, err := json.Marshal(value); err == nil {
			active, _ = strconv.Atoi(activeStr)
			json.Unmarshal(data, &account)
			account.EMail = email
			account.Domain = domain
			account.Provider = provider
			account.Active = active
			account.Scope = scope

			token, err := auth.GenerateToken(&account)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
				return
			}
			account.Token = token.AccessToken
			account.Expired = token.Expire

			if err := (*db).Save(config.ServiceAccountTable, name, account); err != nil {
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
			err = (*db).Delete(config.ServiceAccountTable, account)
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
		tokenObj, tokenOk := c.Get("token")

		if tokenOk {
			tokenInfo = tokenObj.(*token.TokenInfo)
			value, _ := (*db).Load(config.ServiceAccountTable, name)

			if data, err := json.Marshal(value); err == nil {
				json.Unmarshal(data, &account)
				tokenInfo, err = auth.RefreshToken(tokenInfo.AccessToken)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
					account.Token = tokenInfo.AccessToken
					account.Expired = tokenInfo.Expire
					// Save
					if err := (*db).Save(config.ServiceAccountTable, name, account); err != nil {
						c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"token":   tokenInfo,
					})
				}
			}
		}
	}
}

// HandleGetScopeTags process get scope tags
func HandleGetScopeTags() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, config.Scope.Tags())
	}
}
