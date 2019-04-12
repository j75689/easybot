package context

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fatih/structs"
	messagehandler "github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/store"
	"github.com/j75689/easybot/pkg/util"
	"github.com/j75689/easybot/plugin"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/config"
)

// HandleGetConfig process get config file
func HandleGetConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {

		if data, err := (*db).Load(config.MessageHandlerConfigTable, c.Param("id")); err == nil {
			var messageConfig config.MessageHandlerConfig
			b, _ := json.Marshal(data)
			json.Unmarshal(b, &messageConfig)
			c.JSON(200, messageConfig)
		} else {
			c.JSON(200, gin.H{"error": err.Error()})
		}
	}
}

// HandlePostConfig process post config file
func HandlePostConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		if configData, err := c.GetRawData(); err == nil {
			var messageConfig config.MessageHandlerConfig
			if err = json.Unmarshal(configData, &messageConfig); err == nil {
				if err = (*db).Save(config.MessageHandlerConfigTable, messageConfig.ID, messageConfig); err != nil {
					logger.Errorf("[dashboard] Save config [%s] error: %s", messageConfig.ID, err.Error())
				} else {
					logger.Infof("[dashboard] Register config [%s]", messageConfig.ID)
					messagehandler.RegisterConfig(&messageConfig)
					c.JSON(200, gin.H{"message": "success."})
				}
			} else {
				c.JSON(200, gin.H{"error": "invalid config."})
			}

		} else {
			c.JSON(200, gin.H{"error": err.Error()})
		}
	}
}

// HandleDeleteConfig process delete config file
func HandleDeleteConfig(db *store.Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			configID = c.Param("id")
		)
		if data, err := (*db).Load(config.MessageHandlerConfigTable, configID); err == nil {
			var messageConfig config.MessageHandlerConfig
			if b, err := json.Marshal(data); err == nil {
				if err = json.Unmarshal(b, &messageConfig); err != nil {
					logger.Error("[dashboard] ", err.Error())
				} else {
					logger.Infof("[dashboard] Deregister config [%s]", messageConfig.ID)
					if err = messagehandler.DeregisterConfig(&messageConfig); err != nil {
						logger.Error("[dashboard] ", err.Error())
					}
				}
			}

		} else {
			c.JSON(200, gin.H{"error": err.Error()})
			return
		}
		if err := (*db).Delete("config", configID); err != nil {
			logger.Errorf("[dashboard] Delete config [%s] error: %s", configID, err.Error())
			c.JSON(200, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "success."})
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
