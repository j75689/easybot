package context

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/j75689/easybot/model"

	"github.com/fatih/structs"
	"github.com/j75689/easybot/config"
	messagehandler "github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
	"github.com/j75689/easybot/pkg/util"
	"github.com/j75689/easybot/plugin"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/gin-gonic/gin"
)

// HandleGetAllConfigID process get all config id
func HandleGetAllConfigID(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var configIDs = make(map[string][]string)
		if err := (*db).LoadAll(config.MessageHandlerConfigTable, func(id string, value []byte) {
			var messageConfig model.MessageHandlerConfig
			if err := json.Unmarshal(value, &messageConfig); err == nil {
				configIDs[messageConfig.EventType] = append(configIDs[messageConfig.EventType], messageConfig.ConfigID)
			}
		}); err != nil {
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
		}
		c.JSON(200, configIDs)
	}
}

// HandleGetConfig process get config file
func HandleGetConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {

		if data, err := (*db).LoadWithFilter(config.MessageHandlerConfigTable, map[string]interface{}{"configId": c.Param("id")}); err == nil {
			var messageConfig model.MessageHandlerConfig
			json.Unmarshal(data, &messageConfig)
			c.JSON(200, messageConfig)
		} else {
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
		}
	}
}

// HandleCreateConfig process create config file
func HandleCreateConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var configID = c.Param("id")
		if value, err := (*db).LoadWithFilter(config.MessageHandlerConfigTable, map[string]interface{}{"configId": configID}); err == nil {
			logger.Warn("[dashboard] create config already exist:", string(value))
			c.JSON(http.StatusOK, gin.H{"success": false, "error": configID + " already exist"})
			return
		}

		if configData, err := c.GetRawData(); err == nil {
			var messageConfig model.MessageHandlerConfig
			if err = json.Unmarshal(configData, &messageConfig); err == nil {
				if configID != messageConfig.ConfigID {
					c.JSON(http.StatusOK, gin.H{"success": false, "error": "ConfigID not match"})
					return
				}
				if err = (*db).Save(config.MessageHandlerConfigTable, &messageConfig); err != nil {
					logger.Errorf("[dashboard] Save config [%s] error: %s", messageConfig.ConfigID, err.Error())
				} else {
					logger.Infof("[dashboard] Register config [%s]", messageConfig.ConfigID)
					messagehandler.RegisterConfig(&messageConfig)
					c.JSON(200, gin.H{"success": true})
				}
			} else {
				logger.Error("[dashboard] ", err)
				c.JSON(200, gin.H{"success": false, "error": "Invalid config"})
			}

		} else {
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
		}
	}
}

// HandleSaveConfig process save config file
func HandleSaveConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var configID = c.Param("id")
		if configData, err := c.GetRawData(); err == nil {
			var messageConfig model.MessageHandlerConfig
			if err = json.Unmarshal(configData, &messageConfig); err == nil {
				if configID != messageConfig.ConfigID {
					c.JSON(http.StatusOK, gin.H{"success": false, "error": "ConfigID not match"})
					return
				}
				if err = (*db).SaveWithFilter(config.MessageHandlerConfigTable, &messageConfig, map[string]interface{}{"configId": messageConfig.ConfigID}); err != nil {
					logger.Errorf("[dashboard] Save config [%s] error: %s", messageConfig.ConfigID, err.Error())
				} else {
					logger.Infof("[dashboard] Register config [%s]", messageConfig.ConfigID)
					messagehandler.RegisterConfig(&messageConfig)
					c.JSON(200, gin.H{"success": true})
				}
			} else {
				c.JSON(200, gin.H{"success": false, "error": "Invalid config"})
			}

		} else {
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
		}
	}
}

// HandleDeleteConfig process delete config file
func HandleDeleteConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			configID = c.Param("id")
		)
		if data, err := (*db).LoadWithFilter(config.MessageHandlerConfigTable, map[string]interface{}{"configId": configID}); err == nil {
			var messageConfig model.MessageHandlerConfig
			if err = json.Unmarshal(data, &messageConfig); err != nil {
				logger.Error("[dashboard] ", err.Error())
			} else {
				logger.Infof("[dashboard] Deregister config [%s]", messageConfig.ConfigID)
				if err = messagehandler.DeregisterConfig(&messageConfig); err != nil {
					logger.Error("[dashboard] ", err.Error())
				}
			}

		} else {
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
			return
		}
		if err := (*db).DeleteWithFilter(config.MessageHandlerConfigTable, map[string]interface{}{"configId": configID}); err != nil {
			logger.Errorf("[dashboard] Delete config [%s] error: %s", configID, err.Error())
			c.JSON(200, gin.H{"success": false, "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"success": true})
		}
	}
}

// HandleTestRunner process test event handler
func HandleTestRunner() func(*gin.Context) {
	return func(c *gin.Context) {

		defer c.Done()
		postdata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("[dashboard] ", err)
			return
		}
		var arg linebot.Event
		err = json.Unmarshal(postdata, &arg)
		if err != nil {
			logger.Error("[dashboard] ", err)
		}
		logger.Debug("[dashboard] ", structs.Map(arg))
		reply, err := messagehandler.Execute(&arg)
		if err != nil {
			logger.Debug("[dashboard] ", err)
		}

		c.JSON(200, reply)
	}
}

// HandleTestPlugin process test plugin func
func HandleTestPlugin() func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			pluginName = c.Param("plugin")
		)
		defer c.Done()

		type args struct {
			Input     interface{}            `json:"input"`
			Variables map[string]interface{} `json:"variables"`
		}

		postdata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("[dashboard] ", err)
			return
		}

		var arg = args{}
		err = json.Unmarshal(postdata, &arg)
		if err != nil {
			logger.Error("[dashboard] ", err)
		}
		logger.Debug("[dashboard] ", arg)
		var b []byte
		if b, err = json.Marshal(arg.Input); err != nil {
			logger.Error("[dashboard] ", err.Error())
		}
		ParamData := util.ReplaceVariables(string(b), arg.Variables)
		if err = json.Unmarshal([]byte(ParamData), &arg.Input); err != nil {
			logger.Error("[dashboard] ", err.Error())
		}
		v, next, err := plugin.Excute(pluginName, arg.Input, arg.Variables)

		c.JSON(200, gin.H{
			"variables": v,
			"next":      next,
			"error":     err,
		})

	}
}
