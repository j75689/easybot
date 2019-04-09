package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/structs"

	"github.com/j75689/easybot/plugin"

	"github.com/j75689/easybot/config"
	messagehandler "github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/pkg/util"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

func initRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.HTMLRender = gorice.New(rice.MustFindBox("../dashboard/build"))
	// Session
	store := cookie.NewStore([]byte(appName))
	router.Use(sessions.Sessions("session", store))

	// Front-end file
	if _, err := os.Stat("dashboard/build"); err == nil {
		logger.Info("Serve dashboard/build")
		router.Static("/public", "dashboard/build")
	} else {
		logger.Info("Serve rice ", err)
		dist := rice.MustFindBox("../dashboard/build")
		router.StaticFS("/public", dist.HTTPBox())
	}

	// Config CRUD
	router.Any("/config/:id", handleCRUDConfig)

	// Tester
	router.Any("/debug/runner", handleTestRunner)
	router.Any("/debug/plugin/:plugin", handleTestPlugin)

	// Bot
	BotAPI := router.Group("/bot")

	// Line
	LineAPIGroup := BotAPI.Group("/line")
	if handler, err := httphandler.New(channel_secret, channel_token); err == nil {
		if bot, err := handler.NewClient(); err == nil {
			// Reply Hook
			LineAPIGroup.Any("/hook", gin.WrapH(newLineHookHandler(handler, bot)))
			// Push Message
			LineAPIGroup.POST("/push/:userID", handlePushMessage(bot))
			LineAPIGroup.POST("/multicast", handleMulticastMessage(bot))
		} else {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}

	ManagerGroup := router.Group("/")
	ManagerGroup.Use(SessionMiddleware())
	// Dashboard
	{
		// Page API
		ManagerGroup.POST("/login", handleLogin)
		// Pages
		ManagerGroup.GET("/", handleIndexPage)
		ManagerGroup.GET("/dashboard", handleIndexPage)
		ManagerGroup.GET("/login", handleIndexPage)
	}

	return
}

func handleIndexPage(c *gin.Context) {

	if c.Request.URL.Path == "/" {
		c.Redirect(301, context_path+"dashboard")
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func newLineHookHandler(handler *httphandler.WebhookHandler, bot *linebot.Client) *httphandler.WebhookHandler {

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {

		logger.Debug(r)
		for _, event := range events {
			logger.Debug(structs.Map(event))
			if msg, err := messagehandler.Execute(event); msg != nil {
				if err != nil {
					logger.Warn(err)
				}
				msgData, _ := msg.MarshalJSON()
				logger.Debug(string(msgData))
				if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
					logger.Error(err)
				}
			}
		}
	})

	return handler
}

func handlePushMessage(bot *linebot.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			postdata []byte
			err      error
		)
		if postdata, err = ioutil.ReadAll(c.Request.Body); err == nil {
			if _, err = bot.PushMessage(c.Param("userID"), &config.CustomMessage{Msg: string(postdata)}).Do(); err == nil {
				c.JSON(http.StatusOK, map[string]interface{}{
					"success": true,
				})
				return
			}
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func handleMulticastMessage(bot *linebot.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		type MulticastBody struct {
			UserIDs []string      `json:"to"`
			Message []interface{} `json:"messages"`
		}
		var (
			postdata      []byte
			multicastBody MulticastBody
			Messages      []linebot.SendingMessage
			err           error
		)
		if postdata, err = ioutil.ReadAll(c.Request.Body); err == nil {
			if err = json.Unmarshal(postdata, &multicastBody); err == nil {
				for _, data := range multicastBody.Message {
					if msg, err := json.Marshal(data); err == nil {
						Messages = append(Messages, &config.CustomMessage{Msg: string(msg)})
					} else {
						logger.Error("Muticast Cause Error: ", err)
					}
				}

				if _, err = bot.Multicast(multicastBody.UserIDs, Messages...).Do(); err == nil {
					c.JSON(http.StatusOK, map[string]interface{}{
						"success": true,
					})
					return
				}
			}
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"error": err.Error(),
		})

	}
}

func handleCRUDConfig(c *gin.Context) {
	var (
		configID = c.Param("id")
	)

	switch c.Request.Method {
	case "GET":
		if data, err := db.Load("config", configID); err == nil {
			var messageConfig config.MessageHandlerConfig
			b, _ := json.Marshal(data)
			json.Unmarshal(b, &messageConfig)
			c.JSON(200, messageConfig)
		} else {
			c.JSON(200, map[string]string{"error": err.Error()})
		}

	case "POST":
		if configData, err := c.GetRawData(); err == nil {
			var messageConfig config.MessageHandlerConfig
			if err = json.Unmarshal(configData, &messageConfig); err == nil {
				if err = db.Save("config", messageConfig.ID, messageConfig); err != nil {
					logger.Errorf("Save config [%s] error: %s", messageConfig.ID, err.Error())
				} else {
					logger.Infof("Register config [%s]", messageConfig.ID)
					messagehandler.RegisterConfig(&messageConfig)
					c.JSON(200, map[string]string{"message": "success."})
				}
			} else {
				c.JSON(200, map[string]string{"error": "invalid config."})
			}

		} else {
			c.JSON(200, map[string]string{"error": err.Error()})
		}
	case "DELETE":
		if data, err := db.Load("config", configID); err == nil {
			var messageConfig config.MessageHandlerConfig
			if b, err := json.Marshal(data); err == nil {
				if err = json.Unmarshal(b, &messageConfig); err != nil {
					logger.Error(err.Error())
				} else {
					logger.Infof("Deregister config [%s]", messageConfig.ID)
					if err = messagehandler.DeregisterConfig(&messageConfig); err != nil {
						logger.Error(err.Error())
					}
				}
			}

		} else {
			c.JSON(200, map[string]string{"error": err.Error()})
			return
		}
		if err := db.Delete("config", configID); err != nil {
			logger.Errorf("Delete config [%s] error: %s", configID, err.Error())
			c.JSON(200, map[string]string{"error": err.Error()})
		} else {
			c.JSON(200, map[string]string{"message": "success."})
		}
	default:
		c.JSON(405, map[string]string{"error": "Method Not Allowed."})
	}
}

func handleTestRunner(c *gin.Context) {

	defer c.Done()
	postdata, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	var arg linebot.Event
	err = json.Unmarshal(postdata, &arg)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(structs.Map(arg))
	reply, err := messagehandler.Execute(&arg)
	if err != nil {
		logger.Debug(err)
	}

	c.JSON(200, reply)
}

func handleTestPlugin(c *gin.Context) {
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
		logger.Error(err)
		return
	}

	var arg = args{}
	err = json.Unmarshal(postdata, &arg)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(arg)
	var b []byte
	if b, err = json.Marshal(arg.Input); err != nil {
		logger.Error(err.Error())
	}
	ParamData := util.ReplaceVariables(string(b), arg.Variables)
	if err = json.Unmarshal([]byte(ParamData), &arg.Input); err != nil {
		logger.Error(err.Error())
	}
	v, next, err := plugin.Excute(pluginName, arg.Input, arg.Variables)

	c.JSON(200, map[string]interface{}{
		"variables": v,
		"next":      next,
		"error":     err,
	})

}

func handleLogin(c *gin.Context) {
	var (
		user   = c.DefaultPostForm("user", "")
		pass   = c.DefaultPostForm("pass", "")
		result = map[string]interface{}{
			"success": false,
		}
	)

	if user == admin_user && pass == admin_pass {
		session := sessions.Default(c)
		session.Set("login", true)
		session.Save()
		result["success"] = true
	}

	c.JSON(200, result)
}

func SessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		var isLogin = false
		session := sessions.Default(c)
		v := session.Get("login")
		if v != nil {
			isLogin = v.(bool)
		}
		if isLogin || strings.Index(c.Request.URL.Path, "/login") > -1 {
			c.Next()
		} else {
			c.Redirect(http.StatusMovedPermanently, context_path+"login")
			c.Done()
		}

	}

}
