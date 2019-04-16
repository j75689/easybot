package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/model"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

// UserIptableMiddleware handle auth
func UserIptableMiddleware(db *store.Storage, skipper RouteSkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		if skipper != nil {
			if skipper(c) {
				c.Next()
				return
			}
		}
		var pass = true
		if scope, ok := config.Scope.GetScope(c.Request.URL.Path); ok {
			err := (*db).LoadAllWithFilter(config.IpTable, map[string]interface{}{"scope": scope},
				func(key string, value interface{}) {
					var iptable model.Iptables
					data, err := json.Marshal(value)
					if err != nil {
						logger.Warnf("[iptable] key:%v err:%v", key, err)
						return
					}
					err = json.Unmarshal(data, &iptable)
					if err != nil {
						logger.Warnf("[iptable] key:%v err:%v", key, err)
						return
					}

					pass = pass && iptable.Pass(c.ClientIP())
				})
			if err != nil {
				logger.Warn("[iptable] ", err)
				return
			}
			if !pass {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("The IP Address %s is not allowed", c.ClientIP())})
				c.Abort()
			}

		}
	}
}
