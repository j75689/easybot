package context

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/model"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
)

// HandleGetIptables get all iptable
func HandleGetIptables(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {

		var iptables = []*model.Iptable{}

		(*db).LoadAll(config.IpTable, func(id string, value interface{}) {

			var iptable model.Iptable
			data, err := json.Marshal(value)
			if err != nil {
				logger.Warnf("[dashboard] get Iptable id:%v err:%v", id, err)
				return
			}
			err = json.Unmarshal(data, &iptable)
			if err != nil {
				logger.Warnf("[dashboard] get Iptable id:%v err:%v", id, err)
				return
			}
			iptables = append(iptables, &iptable)

		})

		c.JSON(http.StatusOK, iptables)
	}
}

// HandleGetIptable get iptable
func HandleGetIptable(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			id      = c.Param("id")
			iptable model.Iptable
		)
		value, err := (*db).Load(config.IpTable, id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		data, err := json.Marshal(value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		err = json.Unmarshal(data, &iptable)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, iptable)
	}
}

// HandleCreateIptable create new iptable
func HandleCreateIptable(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			iptable model.Iptable
		)
		data, _ := c.GetRawData()
		err := json.Unmarshal(data, &iptable)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		err = (*db).Save(config.IpTable, &iptable)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

// HandleSaveIptable save iptable
func HandleSaveIptable(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			id      = c.Param("id")
			iptable model.Iptable
		)
		data, _ := c.GetRawData()
		err := json.Unmarshal(data, &iptable)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		iptable.ID = id
		err = (*db).Save(config.IpTable, &iptable)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

// HandleDeleteIptable delete iptable
func HandleDeleteIptable(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			id = c.Param("id")
		)

		err := (*db).Delete(config.IpTable, id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
